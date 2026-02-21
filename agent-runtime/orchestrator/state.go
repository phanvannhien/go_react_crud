package orchestrator

type State string

const (
	StatePlan      State = "PLAN"
	StateImplement State = "IMPLEMENT"
	StateValidate  State = "VALIDATE"
	StateFix       State = "FIX"
	StateDone      State = "DONE"
)

type AgentState struct {
	Request        string
	Plan           interface{}
	Implementation interface{}
	Errors         []string
	RetryCount     int
}
