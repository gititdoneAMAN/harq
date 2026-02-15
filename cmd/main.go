package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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

	bot := agent.NewAgent(client)

	if prompt != "" {
        bot.Chat(prompt);
    }

    // Start Interactive Loop
    fmt.Println("\nInteractive Mode Started. Type 'exit' to quit.")
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Print("\n>> ") // Prompt symbol
        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)

        if input == "exit" || input == "quit" {
            fmt.Println("Goodbye!")
            break
        }

        if input == "" {
            continue
        }

        bot.Chat(input)
    }
}

