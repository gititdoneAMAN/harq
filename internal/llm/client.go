package llm

import (
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func NewClient() *openai.Client {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("Env variable API_KEY not found")
	}

	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		panic("Base url is empty.")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))

	return  &client
}