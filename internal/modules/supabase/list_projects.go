package supabase

import (
	"os"

	"github.com/shibaleo/go-mcp-dev/internal/httpclient"
	"github.com/shibaleo/go-mcp-dev/internal/mcp"
)

type ListProjectsTool struct {
	accessToken string
	client      *httpclient.Client
}

func NewListProjectsTool() *ListProjectsTool {
	return &ListProjectsTool{
		accessToken: os.Getenv("SUPABASE_ACCESS_TOKEN"),
		client:      httpclient.New(),
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

	headers := map[string]string{
		"Authorization": "Bearer " + t.accessToken,
	}

	respBody, err := t.client.DoJSON("GET", url, headers, nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}
