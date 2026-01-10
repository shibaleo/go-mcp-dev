package jira

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"

	"github.com/shibaleo/go-mcp-dev/internal/httpclient"
	"github.com/shibaleo/go-mcp-dev/internal/modules"
)

const (
	jiraAPIPath    = "/rest/api/3"
	jiraAPIVersion = "3" // Jira Cloud REST API version
)

var client = httpclient.New()

func getDomain() string {
	return os.Getenv("JIRA_DOMAIN")
}

func getEmail() string {
	return os.Getenv("JIRA_EMAIL")
}

func getAPIToken() string {
	return os.Getenv("JIRA_API_TOKEN")
}

func headers() map[string]string {
	auth := base64.StdEncoding.EncodeToString([]byte(getEmail() + ":" + getAPIToken()))
	return map[string]string{
		"Authorization": "Basic " + auth,
		"Accept":        "application/json",
	}
}

func baseURL() string {
	return fmt.Sprintf("https://%s%s", getDomain(), jiraAPIPath)
}

// Module returns the Jira module definition
func Module() modules.ModuleDefinition {
	return modules.ModuleDefinition{
		Name:        "jira",
		Description: "Jira API - Issue/Project操作（検索、作成、更新、コメント、ワークログ）",
		APIVersion:  jiraAPIVersion,
		TestedAt:    "2026-01-10",
		Tools:       tools,
		Handlers:    handlers,
	}
}

var tools = []modules.Tool{
	{
		Name:        "get_myself",
		Description: "Get information about the current Jira user (myself).",
		InputSchema: modules.InputSchema{
			Type:       "object",
			Properties: map[string]modules.Property{},
		},
	},
	{
		Name:        "list_projects",
		Description: "List all Jira projects accessible to the current user.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"start_at": {
					Type:        "number",
					Description: "Starting index for pagination. Default: 0",
				},
				"max_results": {
					Type:        "number",
					Description: "Maximum results to return. Default: 50",
				},
			},
		},
	},
	{
		Name:        "get_project",
		Description: "Get details of a specific Jira project.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_key": {
					Type:        "string",
					Description: "Project key (e.g., 'PROJ') or ID",
				},
			},
			Required: []string{"project_key"},
		},
	},
	{
		Name:        "search",
		Description: "Search for Jira issues using JQL (Jira Query Language). Example JQL: 'project = PROJ AND status = \"In Progress\"'",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"jql": {
					Type:        "string",
					Description: "JQL query string. Example: 'project = PROJ AND status != Done ORDER BY created DESC'",
				},
				"start_at": {
					Type:        "number",
					Description: "Starting index for pagination. Default: 0",
				},
				"max_results": {
					Type:        "number",
					Description: "Maximum results to return. Default: 50",
				},
				"fields": {
					Type:        "array",
					Description: "Fields to return. Default: summary, status, priority, assignee, created, updated",
				},
			},
			Required: []string{"jql"},
		},
	},
	{
		Name:        "get_issue",
		Description: "Get details of a specific Jira issue by key or ID.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"issue_key": {
					Type:        "string",
					Description: "Issue key (e.g., 'PROJ-123') or ID",
				},
				"fields": {
					Type:        "array",
					Description: "Specific fields to return. If not specified, returns common fields.",
				},
			},
			Required: []string{"issue_key"},
		},
	},
	{
		Name:        "create_issue",
		Description: "Create a new Jira issue.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_key": {
					Type:        "string",
					Description: "Project key (e.g., 'PROJ')",
				},
				"issue_type": {
					Type:        "string",
					Description: "Issue type (e.g., 'Task', 'Bug', 'Story', 'Epic')",
				},
				"summary": {
					Type:        "string",
					Description: "Issue summary/title",
				},
				"description": {
					Type:        "string",
					Description: "Issue description",
				},
				"assignee_account_id": {
					Type:        "string",
					Description: "Assignee's Atlassian account ID",
				},
				"priority": {
					Type:        "string",
					Description: "Priority name (e.g., 'High', 'Medium', 'Low')",
				},
				"labels": {
					Type:        "array",
					Description: "Labels to add to the issue",
				},
				"parent_key": {
					Type:        "string",
					Description: "Parent issue key for subtasks",
				},
			},
			Required: []string{"project_key", "issue_type", "summary"},
		},
	},
	{
		Name:        "update_issue",
		Description: "Update an existing Jira issue.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"issue_key": {
					Type:        "string",
					Description: "Issue key (e.g., 'PROJ-123')",
				},
				"summary": {
					Type:        "string",
					Description: "New summary/title",
				},
				"description": {
					Type:        "string",
					Description: "New description",
				},
				"assignee_account_id": {
					Type:        "string",
					Description: "New assignee's Atlassian account ID",
				},
				"priority": {
					Type:        "string",
					Description: "New priority name",
				},
				"labels": {
					Type:        "array",
					Description: "New labels (replaces existing)",
				},
			},
			Required: []string{"issue_key"},
		},
	},
	{
		Name:        "get_transitions",
		Description: "Get available transitions for an issue. Use this to find valid transition IDs before changing issue status.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"issue_key": {
					Type:        "string",
					Description: "Issue key (e.g., 'PROJ-123')",
				},
			},
			Required: []string{"issue_key"},
		},
	},
	{
		Name:        "transition_issue",
		Description: "Transition an issue to a new status. Use get_transitions first to get valid transition IDs.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"issue_key": {
					Type:        "string",
					Description: "Issue key (e.g., 'PROJ-123')",
				},
				"transition_id": {
					Type:        "string",
					Description: "Transition ID (get from get_transitions)",
				},
				"comment": {
					Type:        "string",
					Description: "Optional comment to add with the transition",
				},
			},
			Required: []string{"issue_key", "transition_id"},
		},
	},
	{
		Name:        "get_comments",
		Description: "Get comments on a Jira issue.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"issue_key": {
					Type:        "string",
					Description: "Issue key (e.g., 'PROJ-123')",
				},
				"start_at": {
					Type:        "number",
					Description: "Starting index for pagination. Default: 0",
				},
				"max_results": {
					Type:        "number",
					Description: "Maximum results to return. Default: 50",
				},
			},
			Required: []string{"issue_key"},
		},
	},
	{
		Name:        "add_comment",
		Description: "Add a comment to a Jira issue.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"issue_key": {
					Type:        "string",
					Description: "Issue key (e.g., 'PROJ-123')",
				},
				"body": {
					Type:        "string",
					Description: "Comment text",
				},
			},
			Required: []string{"issue_key", "body"},
		},
	},
	{
		Name:        "get_worklogs",
		Description: "Get work logs for a Jira issue.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"issue_key": {
					Type:        "string",
					Description: "Issue key (e.g., 'PROJ-123')",
				},
				"start_at": {
					Type:        "number",
					Description: "Starting index for pagination. Default: 0",
				},
				"max_results": {
					Type:        "number",
					Description: "Maximum results to return. Default: 50",
				},
			},
			Required: []string{"issue_key"},
		},
	},
	{
		Name:        "add_worklog",
		Description: "Add a work log to a Jira issue.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"issue_key": {
					Type:        "string",
					Description: "Issue key (e.g., 'PROJ-123')",
				},
				"time_spent_seconds": {
					Type:        "number",
					Description: "Time spent in seconds",
				},
				"started": {
					Type:        "string",
					Description: "Start time in ISO 8601 format (e.g., '2024-01-15T10:00:00.000+0900'). Defaults to now.",
				},
				"comment": {
					Type:        "string",
					Description: "Work log comment",
				},
			},
			Required: []string{"issue_key", "time_spent_seconds"},
		},
	},
}

var handlers = map[string]modules.ToolHandler{
	"get_myself":       getMyself,
	"list_projects":    listProjects,
	"get_project":      getProject,
	"search":           search,
	"get_issue":        getIssue,
	"create_issue":     createIssue,
	"update_issue":     updateIssue,
	"get_transitions":  getTransitions,
	"transition_issue": transitionIssue,
	"get_comments":     getComments,
	"add_comment":      addComment,
	"get_worklogs":     getWorklogs,
	"add_worklog":      addWorklog,
}

// =============================================================================
// User
// =============================================================================

func getMyself(params map[string]interface{}) (string, error) {
	endpoint := baseURL() + "/myself"

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Projects
// =============================================================================

func listProjects(params map[string]interface{}) (string, error) {
	startAt := 0
	if sa, ok := params["start_at"].(float64); ok {
		startAt = int(sa)
	}

	maxResults := 50
	if mr, ok := params["max_results"].(float64); ok {
		maxResults = int(mr)
	}

	endpoint := fmt.Sprintf("%s/project/search?startAt=%d&maxResults=%d", baseURL(), startAt, maxResults)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getProject(params map[string]interface{}) (string, error) {
	projectKey, ok := params["project_key"].(string)
	if !ok {
		return "", fmt.Errorf("project_key must be a string")
	}

	endpoint := fmt.Sprintf("%s/project/%s", baseURL(), projectKey)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Issues
// =============================================================================

func search(params map[string]interface{}) (string, error) {
	jql, ok := params["jql"].(string)
	if !ok {
		return "", fmt.Errorf("jql must be a string")
	}

	query := url.Values{}
	query.Set("jql", jql)

	startAt := 0
	if sa, ok := params["start_at"].(float64); ok {
		startAt = int(sa)
	}
	query.Set("startAt", fmt.Sprintf("%d", startAt))

	maxResults := 50
	if mr, ok := params["max_results"].(float64); ok {
		maxResults = int(mr)
	}
	query.Set("maxResults", fmt.Sprintf("%d", maxResults))

	if fields, ok := params["fields"].([]interface{}); ok && len(fields) > 0 {
		fieldStrs := make([]string, 0, len(fields))
		for _, f := range fields {
			if fs, ok := f.(string); ok {
				fieldStrs = append(fieldStrs, fs)
			}
		}
		if len(fieldStrs) > 0 {
			query.Set("fields", joinStrings(fieldStrs, ","))
		}
	} else {
		query.Set("fields", "summary,status,priority,assignee,created,updated")
	}

	endpoint := fmt.Sprintf("%s/search/jql?%s", baseURL(), query.Encode())

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getIssue(params map[string]interface{}) (string, error) {
	issueKey, ok := params["issue_key"].(string)
	if !ok {
		return "", fmt.Errorf("issue_key must be a string")
	}

	query := url.Values{}
	if fields, ok := params["fields"].([]interface{}); ok && len(fields) > 0 {
		fieldStrs := make([]string, 0, len(fields))
		for _, f := range fields {
			if fs, ok := f.(string); ok {
				fieldStrs = append(fieldStrs, fs)
			}
		}
		if len(fieldStrs) > 0 {
			query.Set("fields", joinStrings(fieldStrs, ","))
		}
	}

	queryStr := ""
	if len(query) > 0 {
		queryStr = "?" + query.Encode()
	}

	endpoint := fmt.Sprintf("%s/issue/%s%s", baseURL(), issueKey, queryStr)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func createIssue(params map[string]interface{}) (string, error) {
	projectKey, ok := params["project_key"].(string)
	if !ok {
		return "", fmt.Errorf("project_key must be a string")
	}

	issueType, ok := params["issue_type"].(string)
	if !ok {
		return "", fmt.Errorf("issue_type must be a string")
	}

	summary, ok := params["summary"].(string)
	if !ok {
		return "", fmt.Errorf("summary must be a string")
	}

	fields := map[string]interface{}{
		"project":   map[string]string{"key": projectKey},
		"issuetype": map[string]string{"name": issueType},
		"summary":   summary,
	}

	if description, ok := params["description"].(string); ok && description != "" {
		fields["description"] = adfDocument(description)
	}

	if assigneeID, ok := params["assignee_account_id"].(string); ok && assigneeID != "" {
		fields["assignee"] = map[string]string{"accountId": assigneeID}
	}

	if priority, ok := params["priority"].(string); ok && priority != "" {
		fields["priority"] = map[string]string{"name": priority}
	}

	if labels, ok := params["labels"].([]interface{}); ok && len(labels) > 0 {
		labelStrs := make([]string, 0, len(labels))
		for _, l := range labels {
			if ls, ok := l.(string); ok {
				labelStrs = append(labelStrs, ls)
			}
		}
		fields["labels"] = labelStrs
	}

	if parentKey, ok := params["parent_key"].(string); ok && parentKey != "" {
		fields["parent"] = map[string]string{"key": parentKey}
	}

	body := map[string]interface{}{"fields": fields}

	endpoint := baseURL() + "/issue"

	respBody, err := client.DoJSON("POST", endpoint, headers(), body)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func updateIssue(params map[string]interface{}) (string, error) {
	issueKey, ok := params["issue_key"].(string)
	if !ok {
		return "", fmt.Errorf("issue_key must be a string")
	}

	fields := make(map[string]interface{})

	if summary, ok := params["summary"].(string); ok {
		fields["summary"] = summary
	}

	if description, ok := params["description"].(string); ok {
		fields["description"] = adfDocument(description)
	}

	if assigneeID, ok := params["assignee_account_id"].(string); ok {
		fields["assignee"] = map[string]string{"accountId": assigneeID}
	}

	if priority, ok := params["priority"].(string); ok {
		fields["priority"] = map[string]string{"name": priority}
	}

	if labels, ok := params["labels"].([]interface{}); ok {
		labelStrs := make([]string, 0, len(labels))
		for _, l := range labels {
			if ls, ok := l.(string); ok {
				labelStrs = append(labelStrs, ls)
			}
		}
		fields["labels"] = labelStrs
	}

	body := map[string]interface{}{"fields": fields}

	endpoint := fmt.Sprintf("%s/issue/%s", baseURL(), issueKey)

	_, err := client.DoJSON("PUT", endpoint, headers(), body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`{"updated": true, "issue_key": "%s"}`, issueKey), nil
}

// =============================================================================
// Transitions
// =============================================================================

func getTransitions(params map[string]interface{}) (string, error) {
	issueKey, ok := params["issue_key"].(string)
	if !ok {
		return "", fmt.Errorf("issue_key must be a string")
	}

	endpoint := fmt.Sprintf("%s/issue/%s/transitions", baseURL(), issueKey)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func transitionIssue(params map[string]interface{}) (string, error) {
	issueKey, ok := params["issue_key"].(string)
	if !ok {
		return "", fmt.Errorf("issue_key must be a string")
	}

	transitionID, ok := params["transition_id"].(string)
	if !ok {
		return "", fmt.Errorf("transition_id must be a string")
	}

	body := map[string]interface{}{
		"transition": map[string]string{"id": transitionID},
	}

	if comment, ok := params["comment"].(string); ok && comment != "" {
		body["update"] = map[string]interface{}{
			"comment": []map[string]interface{}{
				{
					"add": map[string]interface{}{
						"body": adfDocument(comment),
					},
				},
			},
		}
	}

	endpoint := fmt.Sprintf("%s/issue/%s/transitions", baseURL(), issueKey)

	_, err := client.DoJSON("POST", endpoint, headers(), body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`{"transitioned": true, "issue_key": "%s", "transition_id": "%s"}`, issueKey, transitionID), nil
}

// =============================================================================
// Comments
// =============================================================================

func getComments(params map[string]interface{}) (string, error) {
	issueKey, ok := params["issue_key"].(string)
	if !ok {
		return "", fmt.Errorf("issue_key must be a string")
	}

	startAt := 0
	if sa, ok := params["start_at"].(float64); ok {
		startAt = int(sa)
	}

	maxResults := 50
	if mr, ok := params["max_results"].(float64); ok {
		maxResults = int(mr)
	}

	endpoint := fmt.Sprintf("%s/issue/%s/comment?startAt=%d&maxResults=%d", baseURL(), issueKey, startAt, maxResults)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func addComment(params map[string]interface{}) (string, error) {
	issueKey, ok := params["issue_key"].(string)
	if !ok {
		return "", fmt.Errorf("issue_key must be a string")
	}

	body, ok := params["body"].(string)
	if !ok {
		return "", fmt.Errorf("body must be a string")
	}

	payload := map[string]interface{}{
		"body": adfDocument(body),
	}

	endpoint := fmt.Sprintf("%s/issue/%s/comment", baseURL(), issueKey)

	respBody, err := client.DoJSON("POST", endpoint, headers(), payload)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Worklogs
// =============================================================================

func getWorklogs(params map[string]interface{}) (string, error) {
	issueKey, ok := params["issue_key"].(string)
	if !ok {
		return "", fmt.Errorf("issue_key must be a string")
	}

	startAt := 0
	if sa, ok := params["start_at"].(float64); ok {
		startAt = int(sa)
	}

	maxResults := 50
	if mr, ok := params["max_results"].(float64); ok {
		maxResults = int(mr)
	}

	endpoint := fmt.Sprintf("%s/issue/%s/worklog?startAt=%d&maxResults=%d", baseURL(), issueKey, startAt, maxResults)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func addWorklog(params map[string]interface{}) (string, error) {
	issueKey, ok := params["issue_key"].(string)
	if !ok {
		return "", fmt.Errorf("issue_key must be a string")
	}

	timeSpentSeconds, ok := params["time_spent_seconds"].(float64)
	if !ok {
		return "", fmt.Errorf("time_spent_seconds must be a number")
	}

	payload := map[string]interface{}{
		"timeSpentSeconds": int(timeSpentSeconds),
	}

	if started, ok := params["started"].(string); ok && started != "" {
		payload["started"] = started
	}

	if comment, ok := params["comment"].(string); ok && comment != "" {
		payload["comment"] = adfDocument(comment)
	}

	endpoint := fmt.Sprintf("%s/issue/%s/worklog", baseURL(), issueKey)

	respBody, err := client.DoJSON("POST", endpoint, headers(), payload)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Helpers
// =============================================================================

// adfDocument creates an Atlassian Document Format document with a single paragraph
func adfDocument(text string) map[string]interface{} {
	return map[string]interface{}{
		"type":    "doc",
		"version": 1,
		"content": []map[string]interface{}{
			{
				"type": "paragraph",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": text,
					},
				},
			},
		},
	}
}

// joinStrings joins strings with a separator
func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
