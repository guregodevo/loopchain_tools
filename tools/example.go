package tools

import "fmt"

type ExampleTool struct{}

func (e *ExampleTool) Name() string {
	return "Example Tool"
}

func (e *ExampleTool) Description() string {
	return "An example tool to demonstrate the plugin system"
}

func (e *ExampleTool) Call(input string) (string, error) {
	return fmt.Sprintf("Processed: %s", input), nil
}
