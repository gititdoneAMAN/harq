package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gititdoneAMAN/harq/internal/tools"
	"github.com/openai/openai-go/v3"
)

func RunAgent(client *openai.Client, prompt string){
	params := openai.ChatCompletionNewParams{
		Model: os.Getenv("LLM_MODEL"),
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(prompt),
					},
				},
			},
		},
		Tools: tools.GetToolDefinitions(),
	}

	for {
		resp, err := client.Chat.Completions.New(context.Background(),params)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		
		if len(resp.Choices[0].Message.ToolCalls) > 0 {
			fmt.Printf("DEBUG: Tool calls: %+v\n", resp.Choices[0].Message.ToolCalls)
		}

		toolCalls := resp.Choices[0].Message.ToolCalls

		if len(toolCalls) == 0 {
			fmt.Print(resp.Choices[0].Message.Content)
			break
		}

		params.Messages = append(params.Messages, resp.Choices[0].Message.ToParam());

		for _, toolCall := range toolCalls {
			switch toolCall.Function.Name {
            case "read_file":
				var args map[string]any
				err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
				if err!= nil{
					panic("Error unmarshalling arguments")
				}

				location := args["file_path"].(string)

				content := tools.ReadFile(location);

				params.Messages = append(params.Messages, openai.ToolMessage(content, toolCall.ID))
			case "write_file":
				var args map[string]any
				err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
				if err!=nil{
					panic("Unmarshal error")
				}

				location := args["file_path"].(string)
				content := args["content"].(string)

				message := tools.WriteFile(location, content)
				params.Messages = append(params.Messages, openai.ToolMessage(message, toolCall.ID))
			case "run_bash_command":
				var args map[string]any
				err:=json.Unmarshal([]byte(toolCall.Function.Arguments),&args)
				if err!=nil{
					panic("Unmarshall error")
				}

				command := args["command"].(string)

				message := tools.RunBashCommand(command)
				params.Messages = append(params.Messages, openai.ToolMessage(message, toolCall.ID))
			}
		}
	}
}