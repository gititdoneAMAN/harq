package main

import (
	"flag"
	"log"

	"github.com/gititdoneAMAN/harq/internal/agent"
	"github.com/gititdoneAMAN/harq/internal/llm"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	if prompt == "" {
		panic("Prompt must not be empty")
	}

	client := llm.NewClient();

	agent.RunAgent(client,prompt);
}

