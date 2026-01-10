package confluence

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/shibaleo/go-mcp-dev/internal/httpclient"
	"github.com/shibaleo/go-mcp-dev/internal/modules"
)

const (
	confluenceAPIV2 = "/wiki/api/v2"
	confluenceAPIV1 = "/wiki/rest/api"
)

var client = httpclient.New()

func getDomain() string {
	return os.Getenv("CONFLUENCE_DOMAIN")
}

func getEmail() string {
	return os.Getenv("CONFLUENCE_EMAIL")
}

func getAPIToken() string {
	return os.Getenv("CONFLUENCE_API_TOKEN")
}

func headers() map[string]string {
	auth := base64.StdEncoding.EncodeToString([]byte(getEmail() + ":" + getAPIToken()))
	return map[string]string{
		"Authorization": "Basic " + auth,
		"Accept":        "application/json",
	}
}

func baseURLV2() string {
	return fmt.Sprintf("https://%s%s", getDomain(), confluenceAPIV2)
}

func baseURLV1() string {
	return fmt.Sprintf("https://%s%s", getDomain(), confluenceAPIV1)
}

// Module returns the Confluence module definition
func Module() modules.ModuleDefinition {
	return modules.ModuleDefinition{
		Name:        "confluence",
		Description: "Confluence API - Wiki操作（スペース、ページ、検索、コメント、ラベル）",
		Tools:       tools,
		Handlers:    handlers,
	}
}

var tools = []modules.Tool{
	{
		Name:        "list_spaces",
		Description: "List all Confluence spaces accessible to the current user.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"limit": {
					Type:        "number",
					Description: "Maximum results to return. Default: 25",
				},
				"cursor": {
					Type:        "string",
					Description: "Pagination cursor for next page",
				},
			},
		},
	},
	{
		Name:        "get_space",
		Description: "Get details of a specific Confluence space by ID or key.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"space_id_or_key": {
					Type:        "string",
					Description: "Space ID (numeric) or key (e.g., 'MYSPACE')",
				},
			},
			Required: []string{"space_id_or_key"},
		},
	},
	{
		Name:        "get_pages",
		Description: "List pages in a Confluence space.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"space_id": {
					Type:        "string",
					Description: "Space ID (numeric). Use get_space to get ID from key.",
				},
				"limit": {
					Type:        "number",
					Description: "Maximum results to return. Default: 25",
				},
				"cursor": {
					Type:        "string",
					Description: "Pagination cursor for next page",
				},
			},
			Required: []string{"space_id"},
		},
	},
	{
		Name:        "get_page",
		Description: "Get a Confluence page by ID with its content.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"page_id": {
					Type:        "string",
					Description: "Page ID",
				},
				"body_format": {
					Type:        "string",
					Description: "Body format: storage (XHTML) or atlas_doc_format. Default: storage",
				},
			},
			Required: []string{"page_id"},
		},
	},
	{
		Name:        "create_page",
		Description: "Create a new Confluence page.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"space_id": {
					Type:        "string",
					Description: "Space ID (numeric)",
				},
				"title": {
					Type:        "string",
					Description: "Page title",
				},
				"body": {
					Type:        "string",
					Description: "Page body in storage format (XHTML)",
				},
				"parent_id": {
					Type:        "string",
					Description: "Parent page ID for nested pages",
				},
			},
			Required: []string{"space_id", "title", "body"},
		},
	},
	{
		Name:        "update_page",
		Description: "Update an existing Confluence page.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"page_id": {
					Type:        "string",
					Description: "Page ID",
				},
				"title": {
					Type:        "string",
					Description: "New page title",
				},
				"body": {
					Type:        "string",
					Description: "New page body in storage format (XHTML)",
				},
				"version": {
					Type:        "number",
					Description: "Current version number (must be incremented)",
				},
			},
			Required: []string{"page_id", "title", "body", "version"},
		},
	},
	{
		Name:        "delete_page",
		Description: "Delete a Confluence page.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"page_id": {
					Type:        "string",
					Description: "Page ID",
				},
			},
			Required: []string{"page_id"},
		},
	},
	{
		Name:        "search",
		Description: "Search Confluence content using CQL (Confluence Query Language). Example: 'type=page AND space=MYSPACE AND text~\"keyword\"'",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"cql": {
					Type:        "string",
					Description: "CQL query string",
				},
				"limit": {
					Type:        "number",
					Description: "Maximum results to return. Default: 25",
				},
				"start": {
					Type:        "number",
					Description: "Starting index for pagination. Default: 0",
				},
			},
			Required: []string{"cql"},
		},
	},
	{
		Name:        "get_page_comments",
		Description: "Get comments on a Confluence page.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"page_id": {
					Type:        "string",
					Description: "Page ID",
				},
				"limit": {
					Type:        "number",
					Description: "Maximum results to return. Default: 25",
				},
				"cursor": {
					Type:        "string",
					Description: "Pagination cursor for next page",
				},
			},
			Required: []string{"page_id"},
		},
	},
	{
		Name:        "add_page_comment",
		Description: "Add a comment to a Confluence page.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"page_id": {
					Type:        "string",
					Description: "Page ID",
				},
				"body": {
					Type:        "string",
					Description: "Comment body in storage format (XHTML)",
				},
			},
			Required: []string{"page_id", "body"},
		},
	},
	{
		Name:        "get_page_labels",
		Description: "Get labels on a Confluence page.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"page_id": {
					Type:        "string",
					Description: "Page ID",
				},
			},
			Required: []string{"page_id"},
		},
	},
	{
		Name:        "add_page_label",
		Description: "Add a label to a Confluence page.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"page_id": {
					Type:        "string",
					Description: "Page ID",
				},
				"label": {
					Type:        "string",
					Description: "Label name",
				},
			},
			Required: []string{"page_id", "label"},
		},
	},
}

var handlers = map[string]modules.ToolHandler{
	"list_spaces":       listSpaces,
	"get_space":         getSpace,
	"get_pages":         getPages,
	"get_page":          getPage,
	"create_page":       createPage,
	"update_page":       updatePage,
	"delete_page":       deletePage,
	"search":            search,
	"get_page_comments": getPageComments,
	"add_page_comment":  addPageComment,
	"get_page_labels":   getPageLabels,
	"add_page_label":    addPageLabel,
}

// =============================================================================
// Spaces
// =============================================================================

func listSpaces(params map[string]interface{}) (string, error) {
	query := url.Values{}

	limit := 25
	if l, ok := params["limit"].(float64); ok {
		limit = int(l)
	}
	query.Set("limit", fmt.Sprintf("%d", limit))

	if cursor, ok := params["cursor"].(string); ok && cursor != "" {
		query.Set("cursor", cursor)
	}

	endpoint := fmt.Sprintf("%s/spaces?%s", baseURLV2(), query.Encode())

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getSpace(params map[string]interface{}) (string, error) {
	spaceIDOrKey, ok := params["space_id_or_key"].(string)
	if !ok {
		return "", fmt.Errorf("space_id_or_key must be a string")
	}

	// Check if it's numeric (space ID) or a key
	numericRegex := regexp.MustCompile(`^\d+$`)
	var endpoint string

	if numericRegex.MatchString(spaceIDOrKey) {
		// Use V2 API for numeric ID
		endpoint = fmt.Sprintf("%s/spaces/%s", baseURLV2(), spaceIDOrKey)
	} else {
		// Use V1 API for key
		endpoint = fmt.Sprintf("%s/space/%s", baseURLV1(), spaceIDOrKey)
	}

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Pages
// =============================================================================

func getPages(params map[string]interface{}) (string, error) {
	spaceID, ok := params["space_id"].(string)
	if !ok {
		return "", fmt.Errorf("space_id must be a string")
	}

	query := url.Values{}

	limit := 25
	if l, ok := params["limit"].(float64); ok {
		limit = int(l)
	}
	query.Set("limit", fmt.Sprintf("%d", limit))

	if cursor, ok := params["cursor"].(string); ok && cursor != "" {
		query.Set("cursor", cursor)
	}

	endpoint := fmt.Sprintf("%s/spaces/%s/pages?%s", baseURLV2(), spaceID, query.Encode())

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getPage(params map[string]interface{}) (string, error) {
	pageID, ok := params["page_id"].(string)
	if !ok {
		return "", fmt.Errorf("page_id must be a string")
	}

	bodyFormat := "storage"
	if bf, ok := params["body_format"].(string); ok && bf != "" {
		bodyFormat = bf
	}

	endpoint := fmt.Sprintf("%s/pages/%s?body-format=%s", baseURLV2(), pageID, bodyFormat)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func createPage(params map[string]interface{}) (string, error) {
	spaceID, ok := params["space_id"].(string)
	if !ok {
		return "", fmt.Errorf("space_id must be a string")
	}

	title, ok := params["title"].(string)
	if !ok {
		return "", fmt.Errorf("title must be a string")
	}

	body, ok := params["body"].(string)
	if !ok {
		return "", fmt.Errorf("body must be a string")
	}

	payload := map[string]interface{}{
		"spaceId": spaceID,
		"title":   title,
		"status":  "current",
		"body": map[string]interface{}{
			"representation": "storage",
			"value":          body,
		},
	}

	if parentID, ok := params["parent_id"].(string); ok && parentID != "" {
		payload["parentId"] = parentID
	}

	endpoint := baseURLV2() + "/pages"

	respBody, err := client.DoJSON("POST", endpoint, headers(), payload)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func updatePage(params map[string]interface{}) (string, error) {
	pageID, ok := params["page_id"].(string)
	if !ok {
		return "", fmt.Errorf("page_id must be a string")
	}

	title, ok := params["title"].(string)
	if !ok {
		return "", fmt.Errorf("title must be a string")
	}

	body, ok := params["body"].(string)
	if !ok {
		return "", fmt.Errorf("body must be a string")
	}

	version, ok := params["version"].(float64)
	if !ok {
		return "", fmt.Errorf("version must be a number")
	}

	payload := map[string]interface{}{
		"id":     pageID,
		"title":  title,
		"status": "current",
		"body": map[string]interface{}{
			"representation": "storage",
			"value":          body,
		},
		"version": map[string]interface{}{
			"number": int(version),
		},
	}

	endpoint := fmt.Sprintf("%s/pages/%s", baseURLV2(), pageID)

	respBody, err := client.DoJSON("PUT", endpoint, headers(), payload)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func deletePage(params map[string]interface{}) (string, error) {
	pageID, ok := params["page_id"].(string)
	if !ok {
		return "", fmt.Errorf("page_id must be a string")
	}

	endpoint := fmt.Sprintf("%s/pages/%s", baseURLV2(), pageID)

	_, err := client.DoJSON("DELETE", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return `{"deleted": true}`, nil
}

// =============================================================================
// Search (CQL) - uses V1 API
// =============================================================================

func search(params map[string]interface{}) (string, error) {
	cql, ok := params["cql"].(string)
	if !ok {
		return "", fmt.Errorf("cql must be a string")
	}

	query := url.Values{}
	query.Set("cql", cql)

	limit := 25
	if l, ok := params["limit"].(float64); ok {
		limit = int(l)
	}
	query.Set("limit", fmt.Sprintf("%d", limit))

	start := 0
	if s, ok := params["start"].(float64); ok {
		start = int(s)
	}
	query.Set("start", fmt.Sprintf("%d", start))

	endpoint := fmt.Sprintf("%s/search?%s", baseURLV1(), query.Encode())

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Comments
// =============================================================================

func getPageComments(params map[string]interface{}) (string, error) {
	pageID, ok := params["page_id"].(string)
	if !ok {
		return "", fmt.Errorf("page_id must be a string")
	}

	query := url.Values{}

	limit := 25
	if l, ok := params["limit"].(float64); ok {
		limit = int(l)
	}
	query.Set("limit", fmt.Sprintf("%d", limit))

	if cursor, ok := params["cursor"].(string); ok && cursor != "" {
		query.Set("cursor", cursor)
	}

	endpoint := fmt.Sprintf("%s/pages/%s/footer-comments?%s", baseURLV2(), pageID, query.Encode())

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func addPageComment(params map[string]interface{}) (string, error) {
	pageID, ok := params["page_id"].(string)
	if !ok {
		return "", fmt.Errorf("page_id must be a string")
	}

	body, ok := params["body"].(string)
	if !ok {
		return "", fmt.Errorf("body must be a string")
	}

	payload := map[string]interface{}{
		"body": map[string]interface{}{
			"representation": "storage",
			"value":          body,
		},
	}

	endpoint := fmt.Sprintf("%s/pages/%s/footer-comments", baseURLV2(), pageID)

	respBody, err := client.DoJSON("POST", endpoint, headers(), payload)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Labels
// =============================================================================

func getPageLabels(params map[string]interface{}) (string, error) {
	pageID, ok := params["page_id"].(string)
	if !ok {
		return "", fmt.Errorf("page_id must be a string")
	}

	endpoint := fmt.Sprintf("%s/pages/%s/labels", baseURLV2(), pageID)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func addPageLabel(params map[string]interface{}) (string, error) {
	pageID, ok := params["page_id"].(string)
	if !ok {
		return "", fmt.Errorf("page_id must be a string")
	}

	label, ok := params["label"].(string)
	if !ok {
		return "", fmt.Errorf("label must be a string")
	}

	payload := map[string]interface{}{
		"name": label,
	}

	endpoint := fmt.Sprintf("%s/pages/%s/labels", baseURLV2(), pageID)

	respBody, err := client.DoJSON("POST", endpoint, headers(), payload)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}
