package main

import (
	"context"
	"fmt"
)

// ExampleTool is exported with a capital letter so it's accessible in the plugin.
type ExampleTool struct{ ApiKey string }

// Name returns the name of the tool
func (e *ExampleTool) Name() string {
	return "Example Tool"
}

// Description returns the description of the tool
func (e *ExampleTool) Description() string {
	return "An example tool to demonstrate the plugin system"
}

// Call processes the input string and returns a response
func (e *ExampleTool) Call(ctxt context.Context, input string) (string, error) {
	return fmt.Sprintf("Processed: %s with $s", input, e.ApiKey), nil
}

// Export the tool so it can be accessed via the plugin system
var Tool ExampleTool
