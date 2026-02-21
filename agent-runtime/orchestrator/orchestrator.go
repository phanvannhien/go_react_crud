package orchestrator

import (
	"log"

	"agent-runtime/cache"
	"agent-runtime/diff"
	"agent-runtime/gitutil"
	"agent-runtime/llm"
	"agent-runtime/util"
	"agent-runtime/validator"
)

func Run(request string, client *llm.Client) error {

	hash := util.Hash(request)

	// CACHE CHECK
	if cached, _ := cache.Get(hash); cached != "" {
		log.Println("Cache hit")
		return nil
	}

	state := &AgentState{
		Request:    request,
		RetryCount: 0,
	}

	// PLAN
	plan, err := client.GeneratePlan(request)
	if err != nil {
		return err
	}
	state.Plan = plan

	// BUILD DIFF CONTEXT
	context, err := diff.BuildContext(plan.FilesToModify)
	if err != nil {
		return err
	}

	// IMPLEMENT
	impl, err := client.GenerateImplementation(request, context)
	if err != nil {
		return err
	}
	state.Implementation = impl

	// APPLY CHANGES
	err = gitutil.ApplyChanges(impl.Files)
	if err != nil {
		return err
	}

	// VALIDATE
	errors := validator.RunAll()
	state.Errors = errors

	for len(state.Errors) > 0 && state.RetryCount < 2 {

		state.RetryCount++

		log.Println("Validation failed. Attempt fix...")

		impl, err = client.FixImplementation(state.Errors)
		if err != nil {
			return err
		}

		err = gitutil.ApplyChanges(impl.Files)
		if err != nil {
			return err
		}

		state.Errors = validator.RunAll()
	}

	if len(state.Errors) > 0 {
		return err
	}

	cache.Set(hash, "done")

	return nil
}
