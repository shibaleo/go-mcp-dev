package supabase

import (
	"fmt"
	"os"

	"github.com/shibaleo/go-mcp-dev/internal/httpclient"
	"github.com/shibaleo/go-mcp-dev/internal/modules"
)

var client = httpclient.New()

func getAccessToken() string {
	return os.Getenv("SUPABASE_ACCESS_TOKEN")
}

// Module returns the Supabase module definition
func Module() modules.ModuleDefinition {
	return modules.ModuleDefinition{
		Name:        "supabase",
		Description: "Supabase Management API - プロジェクト管理、SQL実行",
		Tools: []modules.Tool{
			{
				Name:        "list_projects",
				Description: "List all Supabase projects accessible with the current access token",
				InputSchema: modules.InputSchema{
					Type:       "object",
					Properties: map[string]modules.Property{},
				},
			},
			{
				Name:        "run_query",
				Description: "Execute a SQL query against a Supabase project database using the Management API",
				InputSchema: modules.InputSchema{
					Type: "object",
					Properties: map[string]modules.Property{
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
			},
		},
		Handlers: map[string]modules.ToolHandler{
			"list_projects": listProjects,
			"run_query":     runQuery,
		},
	}
}

func listProjects(params map[string]interface{}) (string, error) {
	url := "https://api.supabase.com/v1/projects"

	headers := map[string]string{
		"Authorization": "Bearer " + getAccessToken(),
	}

	respBody, err := client.DoJSON("GET", url, headers, nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func runQuery(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	query, ok := params["query"].(string)
	if !ok {
		return "", fmt.Errorf("query must be a string")
	}

	url := fmt.Sprintf("https://api.supabase.com/v1/projects/%s/database/query", projectRef)

	headers := map[string]string{
		"Authorization": "Bearer " + getAccessToken(),
	}

	payload := map[string]string{"query": query}

	respBody, err := client.DoJSON("POST", url, headers, payload)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}
