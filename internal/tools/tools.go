package tools

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func ReadFile(filePath string) string {
	println("-----------------TOOL-CALL------------------------")
	println("Function called: read_file")
	if filePath == "" {
		return "File path is empty"
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "Error reading file"
	}

	return string(content)
}

func WriteFile(filePath string, content string) string {
	println("-----------------TOOL-CALL------------------------")
	println("Function called: write_file")
	if filePath == "" {
		return "File path is empty"
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "Error writing file"
	}

	return "File written successfully"
}

func RunBashCommand(command string) string {
	println("-----------------TOOL-CALL------------------------")
	println("Function called: run_bash_command")
	if command == "" {
		return "Command is empty"
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error running command: %v\nOutput: %s", err, string(output))
	}

	return string(output)
}