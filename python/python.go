package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// PythonExecutionTool represents the tool that sends the code to the Cloud Function
type PythonExecutionTool struct {
	EndpointUrl string
}

// Define the request and response structure for Cloud Function
type CodeExecutionRequest struct {
	Code string `json:"code"`
}

type CodeExecutionResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

// Name returns the name of the tool
func (t *PythonExecutionTool) Name() string {
	return "PythonExecutionTool"
}

// Description returns a short description of the tool
func (t *PythonExecutionTool) Description() string {
	return "Executes Python code via a Cloud Function and returns the result or error."
}

func cleanPython(code string) string {
	// Remove the "```python" marker
	code = strings.Replace(code, "```python", "", -1)
	// Remove the ending "```" marker
	code = strings.Replace(code, "```", "", -1)
	// Trim any leading or trailing whitespace
	return strings.TrimSpace(code)
}

// Call executes the provided Python code by sending it to the Cloud Function
func (t *PythonExecutionTool) Call(ctx context.Context, code string) (string, error) {
	log.Println("[INFO] Preparing to execute Python code via Cloud Function")

	// Prepare the request payload
	requestBody := CodeExecutionRequest{
		Code: cleanPython(code),
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal request: %v\n", err)
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	log.Printf("[INFO] Sending request to Cloud Function at URL: %s\n", t.EndpointUrl)
	log.Printf("[DEBUG] Request payload: %s\n", bodyBytes)

	// Send HTTP POST request to the Cloud Function
	resp, err := http.Post(t.EndpointUrl, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("[ERROR] Failed to execute request: %v\n", err)
		return "", fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("[INFO] Received response with status code: %d\n", resp.StatusCode)

	// Read the response from the Cloud Function
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read response: %v\n", err)
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	log.Printf("[DEBUG] Response body: %s\n", body)

	// Parse the response
	var executionResponse CodeExecutionResponse
	err = json.Unmarshal(body, &executionResponse)
	if err != nil {
		log.Printf("[ERROR] Failed to unmarshal response: %v\n", err)
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Check for any execution error
	if executionResponse.Error != "" {
		log.Printf("[ERROR] Execution error: %s\n", executionResponse.Error)
		return executionResponse.Error, fmt.Errorf("execution error: %v", executionResponse.Error)
	}

	log.Println("[INFO] Python code executed successfully")
	log.Printf("[DEBUG] Execution result: %s\n", executionResponse.Result)

	// Return the result of the Python code execution
	return executionResponse.Result, nil
}

var Tool PythonExecutionTool
