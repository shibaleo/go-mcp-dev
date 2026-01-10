package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/shibaleo/go-mcp-dev/internal/mcp"
)

type RunQueryTool struct {
	accessToken string
	httpClient  *http.Client
}

func NewRunQueryTool() *RunQueryTool {
	return &RunQueryTool{
		accessToken: os.Getenv("SUPABASE_ACCESS_TOKEN"),
		httpClient:  &http.Client{},
	}
}

func (t *RunQueryTool) Definition() mcp.Tool {
	return mcp.Tool{
		Name:        "supabase_run_query",
		Description: "Execute a SQL query against a Supabase project database using the Management API",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"project_ref": {
					Type:        "string",
					Description: "The Supabase project reference ID",
				},
				"query": {
					Type:        "string",
					Description: "The SQL query to execute",
				},
			},
			Required: []string{"project_ref", "query"},
		},
	}
}

func (t *RunQueryTool) Execute(args map[string]interface{}) (string, error) {
	projectRef, ok := args["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("query must be a string")
	}

	// Supabase Management API endpoint
	url := fmt.Sprintf("https://api.supabase.com/v1/projects/%s/database/query", projectRef)

	payload := map[string]string{"query": query}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+t.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// 200 OK and 201 Created are both success
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// Pretty print JSON response
	var result interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return string(respBody), nil
	}

	prettyJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return string(respBody), nil
	}

	return string(prettyJSON), nil
}
