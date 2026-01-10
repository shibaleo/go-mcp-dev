package supabase

import (
	"fmt"
	"os"

	"github.com/shibaleo/go-mcp-dev/internal/httpclient"
	"github.com/shibaleo/go-mcp-dev/internal/mcp"
)

type RunQueryTool struct {
	accessToken string
	client      *httpclient.Client
}

func NewRunQueryTool() *RunQueryTool {
	return &RunQueryTool{
		accessToken: os.Getenv("SUPABASE_ACCESS_TOKEN"),
		client:      httpclient.New(),
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

	url := fmt.Sprintf("https://api.supabase.com/v1/projects/%s/database/query", projectRef)

	headers := map[string]string{
		"Authorization": "Bearer " + t.accessToken,
	}

	payload := map[string]string{"query": query}

	respBody, err := t.client.DoJSON("POST", url, headers, payload)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}
