package main

import (
	"log"
	"os"

	"agent-runtime/llm"
	"agent-runtime/orchestrator"
)

func main() {

	request := os.Args[1]

	client := llm.NewClient(
		os.Getenv("LLM_ENDPOINT"),
		os.Getenv("LLM_API_KEY"),
	)

	err := orchestrator.Run(request, client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Agent completed successfully")
}
