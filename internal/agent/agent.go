package agent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gititdoneAMAN/harq/internal/prompts"
	"github.com/gititdoneAMAN/harq/internal/tools"
	"github.com/openai/openai-go/v3"
)

type Agent struct {
	Client *openai.Client
	History []openai.ChatCompletionMessageParamUnion
}

func NewAgent(client *openai.Client) *Agent {
	return &Agent{
		Client: client,
		History: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompts.GetSystemPrompt()),
		},
	}
}

func (a *Agent)Chat(userPropmt string){
	a.History = append(a.History, openai.UserMessage(userPropmt));

	fmt.Printf("Thinking.....\n")

	for{
		params := openai.ChatCompletionNewParams{
			Model: os.Getenv("LLM_MODEL"),
			Messages: a.History,
			Tools: tools.GetToolDefinitions(),
		}

		ctx:= context.Background()
		resp, err := a.Client.Chat.Completions.New(ctx, params)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		msg:= resp.Choices[0].Message
		toolCalls := msg.ToolCalls
		a.History = append(a.History, msg.ToParam())

		if len(toolCalls) == 0 {
			fmt.Print(resp.Choices[0].Message.Content)
			break
		}

		// if len(toolCalls) > 0 {
		// 	fmt.Printf("DEBUG: Tool calls: %+v\n", resp.Choices[0].Message.ToolCalls)
		// }

		 for _, toolCall := range toolCalls {
            var result string
            var args map[string]any

            if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
                result = fmt.Sprintf("Error parsing arguments: %v", err)
            } else {
                switch toolCall.Function.Name {
                case "read_file":
                    result = tools.ReadFile(args["file_path"].(string))

                case "write_file":
                    path := args["file_path"].(string)
                    content := args["content"].(string)
                    if askUserPermission(fmt.Sprintf("create/write file '%s'", path)) {
                        result = tools.WriteFile(path, content)
                    } else {
                        result = "User denied file write permission."
                    }

                case "run_bash_command":
                    cmd := args["command"].(string)
                    if askUserPermission(fmt.Sprintf("run command '%s'", cmd)) {
                        result = tools.RunBashCommand(cmd)
                    } else {
                        result = "User denied command execution."
                    }
                default:
                    result = fmt.Sprintf("Tool %s not implemented", toolCall.Function.Name)
                }
            }

            a.History = append(a.History, openai.ToolMessage(result, toolCall.ID))
        }
    }
}

func askUserPermission(action string) bool {
    fmt.Printf("⚠️  Agent wants to %s. Allow? [y/N]: ", action)
    reader := bufio.NewReader(os.Stdin)
    response, _ := reader.ReadString('\n')
    response = strings.TrimSpace(response)
    return strings.ToLower(response) == "y"
}