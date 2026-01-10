package supabase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/shibaleo/go-mcp-dev/internal/mcp"
)

type ListProjectsTool struct {
	accessToken string
	httpClient  *http.Client
}

func NewListProjectsTool() *ListProjectsTool {
	return &ListProjectsTool{
		accessToken: os.Getenv("SUPABASE_ACCESS_TOKEN"),
		httpClient:  &http.Client{},
	}
}

func (t *ListProjectsTool) Definition() mcp.Tool {
	return mcp.Tool{
		Name:        "supabase_list_projects",
		Description: "List all Supabase projects accessible with the current access token",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
		},
	}
}

func (t *ListProjectsTool) Execute(args map[string]interface{}) (string, error) {
	url := "https://api.supabase.com/v1/projects"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+t.accessToken)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
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
