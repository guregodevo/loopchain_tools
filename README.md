# loopchain_tools example

## Plugin System for Go Agent Tools

This repository implements dynamic tool integration for the Go-based agent executor using Go's plugin system. This allows you to easily add and use tools dynamically in your agent executor without rebuilding the entire application.

### Setup Instructions

*Prerequisites*

Make sure you have the following installed:

    * Go (version 1.12 or higher, to support plugins)
    * Make
    * Access to the GitHub repository for downloading plugins

*Build and Push a New Tool*

You can add a tool by implementing it in Go and building it as a .so file using Go's plugin system. Once built, you can push the plugin to GitHub Releases to make it accessible for dynamic loading.

*Create Your Tool:*

Create a new Go tool in the tools package, implementing the Tool interface. Here's an example:

```go
    go

package tools

import "fmt"

// ExampleTool implements the Tool interface
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
```

*Build the Plugin*

Use the provided Makefile to build the tool as a plugin.

```bash

make build_plugin PLUGIN_NAME=exampletool

```

This will produce a .so file for your plugin (e.g., exampletool.so).

Push the Plugin to GitHub Releases:

After building the plugin, you can upload it to your GitHub repository's releases section.

First, ensure you have tagged a release using git tag and git push --tags.
    Then, use the GitHub web interface or CLI tool (gh release create) to upload the .so file as part of the release assets.

```bash
    make release_plugin PLUGIN_NAME=exampletool
```

*Using the Plugin in the Agent Executor*

The agent executor will dynamically download and load the plugin when it runs.

    Download and Load the Plugin:

    The agent executor will automatically download the required plugin from GitHub and load it using the plugin package.

    Here's an example usage:

```go

    agentTools := []tools.Tool{}

    // Dynamically load the exampletool plugin
    pluginPath := "./exampletool.so"
    plug, err := plugin.Open(pluginPath)
    if err != nil {
        log.Fatalf("Failed to load plugin: %v", err)
    }

    sym, err := plug.Lookup("ExampleTool")
    if err != nil {
        log.Fatalf("Failed to find symbol in plugin: %v", err)
    }

    exampleTool := sym.(tools.Tool)
    agentTools = append(agentTools, exampleTool)
```

*Run the Agent Executor:*

After loading the tools, you can run the agent executor as normal.

*Makefile Commands*

The provided Makefile simplifies building and releasing plugins:

    make build_plugin PLUGIN_NAME=yourtool: Builds a Go plugin as a .so file.
    make release_plugin PLUGIN_NAME=yourtool: Pushes the plugin to the GitHub release for the project.

### Example Usage

```bash

make build_plugin PLUGIN_NAME=exampletool
make release_plugin PLUGIN_NAME=exampletool
```

After building and releasing the plugin, you can dynamically load it into the agent executor as shown above.
