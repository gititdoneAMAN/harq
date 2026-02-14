package tools

import "github.com/openai/openai-go/v3"

func GetToolDefinitions() []openai.ChatCompletionToolUnionParam {
    return []openai.ChatCompletionToolUnionParam{
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name: "read_file",
			Description: openai.String("Read content of a file specified at the file path"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"file_path": map[string]string{
						"type": "string",
					},
				},
				"required": []string{"file_path"},
			},
		}),
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name: "write_file",
			Description: openai.String("Write content to a file specified at the file path"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"file_path": map[string]string{
						"type": "string",
					},
					"content": map[string]string{
						"type": "string",
					},
				},
				"required": []string{"file_path", "content"},
			},
		}),
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name: "run_bash_command",
			Description: openai.String("Run a bash command"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"command": map[string]string{
						"type": "string",
					},
				},
				"required": []string{"command"},
			},
		}),
	}
}