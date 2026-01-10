package github

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/shibaleo/go-mcp-dev/internal/httpclient"
	"github.com/shibaleo/go-mcp-dev/internal/modules"
)

const (
	githubAPIBase   = "https://api.github.com"
	githubAPIVersion = "2022-11-28"
)

var client = httpclient.New()

func getToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

func headers() map[string]string {
	return map[string]string{
		"Authorization":        "Bearer " + getToken(),
		"Accept":               "application/vnd.github+json",
		"X-GitHub-Api-Version": githubAPIVersion,
	}
}

// Module returns the GitHub module definition
func Module() modules.ModuleDefinition {
	return modules.ModuleDefinition{
		Name:        "github",
		Description: "GitHub API - リポジトリ、Issue、PR、Actions、検索",
		Tools:       tools,
		Handlers:    handlers,
	}
}

var tools = []modules.Tool{
	// User
	{
		Name:        "get_user",
		Description: "Get information about the authenticated GitHub user.",
		InputSchema: modules.InputSchema{
			Type:       "object",
			Properties: map[string]modules.Property{},
		},
	},
	// Repositories
	{
		Name:        "list_repos",
		Description: "List repositories for the authenticated user.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"type": {
					Type:        "string",
					Description: "Type of repositories (all, owner, public, private). Default: owner",
				},
				"sort": {
					Type:        "string",
					Description: "Sort by (created, updated, pushed, full_name). Default: updated",
				},
				"per_page": {
					Type:        "number",
					Description: "Results per page. Default: 30",
				},
				"page": {
					Type:        "number",
					Description: "Page number. Default: 1",
				},
			},
		},
	},
	{
		Name:        "get_repo",
		Description: "Get details of a specific repository.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
			},
			Required: []string{"owner", "repo"},
		},
	},
	{
		Name:        "list_branches",
		Description: "List branches in a repository.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"per_page": {
					Type:        "number",
					Description: "Results per page. Default: 30",
				},
			},
			Required: []string{"owner", "repo"},
		},
	},
	{
		Name:        "list_commits",
		Description: "List commits in a repository.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"sha": {
					Type:        "string",
					Description: "Branch name or commit SHA to filter by",
				},
				"per_page": {
					Type:        "number",
					Description: "Results per page. Default: 30",
				},
				"page": {
					Type:        "number",
					Description: "Page number. Default: 1",
				},
			},
			Required: []string{"owner", "repo"},
		},
	},
	{
		Name:        "get_file_content",
		Description: "Get the content of a file in a repository.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"path": {
					Type:        "string",
					Description: "File path",
				},
				"ref": {
					Type:        "string",
					Description: "Branch name or commit SHA",
				},
			},
			Required: []string{"owner", "repo", "path"},
		},
	},
	// Issues
	{
		Name:        "list_issues",
		Description: "List issues in a repository.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"state": {
					Type:        "string",
					Description: "Issue state (open, closed, all). Default: open",
				},
				"per_page": {
					Type:        "number",
					Description: "Results per page. Default: 30",
				},
				"page": {
					Type:        "number",
					Description: "Page number. Default: 1",
				},
			},
			Required: []string{"owner", "repo"},
		},
	},
	{
		Name:        "get_issue",
		Description: "Get details of a specific issue.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"issue_number": {
					Type:        "number",
					Description: "Issue number",
				},
			},
			Required: []string{"owner", "repo", "issue_number"},
		},
	},
	{
		Name:        "create_issue",
		Description: "Create a new issue in a repository.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"title": {
					Type:        "string",
					Description: "Issue title",
				},
				"body": {
					Type:        "string",
					Description: "Issue body",
				},
				"labels": {
					Type:        "array",
					Description: "Labels to assign",
				},
				"assignees": {
					Type:        "array",
					Description: "Users to assign",
				},
			},
			Required: []string{"owner", "repo", "title"},
		},
	},
	{
		Name:        "update_issue",
		Description: "Update an existing issue.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"issue_number": {
					Type:        "number",
					Description: "Issue number",
				},
				"title": {
					Type:        "string",
					Description: "New title",
				},
				"body": {
					Type:        "string",
					Description: "New body",
				},
				"state": {
					Type:        "string",
					Description: "New state (open, closed)",
				},
				"labels": {
					Type:        "array",
					Description: "Labels to set",
				},
				"assignees": {
					Type:        "array",
					Description: "Users to assign",
				},
			},
			Required: []string{"owner", "repo", "issue_number"},
		},
	},
	{
		Name:        "add_issue_comment",
		Description: "Add a comment to an issue.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"issue_number": {
					Type:        "number",
					Description: "Issue number",
				},
				"body": {
					Type:        "string",
					Description: "Comment body",
				},
			},
			Required: []string{"owner", "repo", "issue_number", "body"},
		},
	},
	// Pull Requests
	{
		Name:        "list_prs",
		Description: "List pull requests in a repository.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"state": {
					Type:        "string",
					Description: "PR state (open, closed, all). Default: open",
				},
				"per_page": {
					Type:        "number",
					Description: "Results per page. Default: 30",
				},
				"page": {
					Type:        "number",
					Description: "Page number. Default: 1",
				},
			},
			Required: []string{"owner", "repo"},
		},
	},
	{
		Name:        "get_pr",
		Description: "Get details of a specific pull request.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"pr_number": {
					Type:        "number",
					Description: "PR number",
				},
			},
			Required: []string{"owner", "repo", "pr_number"},
		},
	},
	{
		Name:        "create_pr",
		Description: "Create a new pull request.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"title": {
					Type:        "string",
					Description: "PR title",
				},
				"head": {
					Type:        "string",
					Description: "Branch with changes",
				},
				"base": {
					Type:        "string",
					Description: "Branch to merge into",
				},
				"body": {
					Type:        "string",
					Description: "PR description",
				},
				"draft": {
					Type:        "boolean",
					Description: "Create as draft PR",
				},
			},
			Required: []string{"owner", "repo", "title", "head", "base"},
		},
	},
	{
		Name:        "list_pr_commits",
		Description: "List commits in a pull request.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"pr_number": {
					Type:        "number",
					Description: "PR number",
				},
				"per_page": {
					Type:        "number",
					Description: "Results per page. Default: 30",
				},
			},
			Required: []string{"owner", "repo", "pr_number"},
		},
	},
	{
		Name:        "list_pr_files",
		Description: "List files changed in a pull request.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"pr_number": {
					Type:        "number",
					Description: "PR number",
				},
				"per_page": {
					Type:        "number",
					Description: "Results per page. Default: 30",
				},
			},
			Required: []string{"owner", "repo", "pr_number"},
		},
	},
	{
		Name:        "list_pr_reviews",
		Description: "List reviews on a pull request.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"pr_number": {
					Type:        "number",
					Description: "PR number",
				},
			},
			Required: []string{"owner", "repo", "pr_number"},
		},
	},
	// Search
	{
		Name:        "search_repos",
		Description: "Search for repositories.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"query": {
					Type:        "string",
					Description: "Search query",
				},
				"sort": {
					Type:        "string",
					Description: "Sort by (stars, forks, help-wanted-issues, updated)",
				},
				"per_page": {
					Type:        "number",
					Description: "Results per page. Default: 30",
				},
				"page": {
					Type:        "number",
					Description: "Page number. Default: 1",
				},
			},
			Required: []string{"query"},
		},
	},
	{
		Name:        "search_code",
		Description: "Search for code across repositories.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"query": {
					Type:        "string",
					Description: "Search query (e.g., 'addClass in:file language:js repo:jquery/jquery')",
				},
				"per_page": {
					Type:        "number",
					Description: "Results per page. Default: 30",
				},
				"page": {
					Type:        "number",
					Description: "Page number. Default: 1",
				},
			},
			Required: []string{"query"},
		},
	},
	{
		Name:        "search_issues",
		Description: "Search for issues and pull requests.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"query": {
					Type:        "string",
					Description: "Search query (e.g., 'repo:owner/repo is:open is:issue')",
				},
				"sort": {
					Type:        "string",
					Description: "Sort by (comments, reactions, created, updated)",
				},
				"per_page": {
					Type:        "number",
					Description: "Results per page. Default: 30",
				},
				"page": {
					Type:        "number",
					Description: "Page number. Default: 1",
				},
			},
			Required: []string{"query"},
		},
	},
	// Actions
	{
		Name:        "list_workflows",
		Description: "List workflows in a repository.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"per_page": {
					Type:        "number",
					Description: "Results per page. Default: 30",
				},
			},
			Required: []string{"owner", "repo"},
		},
	},
	{
		Name:        "list_workflow_runs",
		Description: "List workflow runs in a repository.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"workflow_id": {
					Type:        "string",
					Description: "Workflow ID or file name to filter by",
				},
				"status": {
					Type:        "string",
					Description: "Filter by status (queued, in_progress, completed)",
				},
				"per_page": {
					Type:        "number",
					Description: "Results per page. Default: 30",
				},
			},
			Required: []string{"owner", "repo"},
		},
	},
	{
		Name:        "get_workflow_run",
		Description: "Get details of a specific workflow run.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"run_id": {
					Type:        "number",
					Description: "Workflow run ID",
				},
			},
			Required: []string{"owner", "repo", "run_id"},
		},
	},
}

var handlers = map[string]modules.ToolHandler{
	"get_user":           getUser,
	"list_repos":         listRepos,
	"get_repo":           getRepo,
	"list_branches":      listBranches,
	"list_commits":       listCommits,
	"get_file_content":   getFileContent,
	"list_issues":        listIssues,
	"get_issue":          getIssue,
	"create_issue":       createIssue,
	"update_issue":       updateIssue,
	"add_issue_comment":  addIssueComment,
	"list_prs":           listPRs,
	"get_pr":             getPR,
	"create_pr":          createPR,
	"list_pr_commits":    listPRCommits,
	"list_pr_files":      listPRFiles,
	"list_pr_reviews":    listPRReviews,
	"search_repos":       searchRepos,
	"search_code":        searchCode,
	"search_issues":      searchIssues,
	"list_workflows":     listWorkflows,
	"list_workflow_runs": listWorkflowRuns,
	"get_workflow_run":   getWorkflowRun,
}

// =============================================================================
// User
// =============================================================================

func getUser(params map[string]interface{}) (string, error) {
	endpoint := githubAPIBase + "/user"

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Repositories
// =============================================================================

func listRepos(params map[string]interface{}) (string, error) {
	query := url.Values{}

	if t, ok := params["type"].(string); ok && t != "" {
		query.Set("type", t)
	} else {
		query.Set("type", "owner")
	}

	if sort, ok := params["sort"].(string); ok && sort != "" {
		query.Set("sort", sort)
	} else {
		query.Set("sort", "updated")
	}

	perPage := 30
	if pp, ok := params["per_page"].(float64); ok {
		perPage = int(pp)
	}
	query.Set("per_page", fmt.Sprintf("%d", perPage))

	page := 1
	if p, ok := params["page"].(float64); ok {
		page = int(p)
	}
	query.Set("page", fmt.Sprintf("%d", page))

	endpoint := fmt.Sprintf("%s/user/repos?%s", githubAPIBase, query.Encode())

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getRepo(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s", githubAPIBase, owner, repo)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func listBranches(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	perPage := 30
	if pp, ok := params["per_page"].(float64); ok {
		perPage = int(pp)
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/branches?per_page=%d", githubAPIBase, owner, repo, perPage)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func listCommits(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	query := url.Values{}

	if sha, ok := params["sha"].(string); ok && sha != "" {
		query.Set("sha", sha)
	}

	perPage := 30
	if pp, ok := params["per_page"].(float64); ok {
		perPage = int(pp)
	}
	query.Set("per_page", fmt.Sprintf("%d", perPage))

	page := 1
	if p, ok := params["page"].(float64); ok {
		page = int(p)
	}
	query.Set("page", fmt.Sprintf("%d", page))

	endpoint := fmt.Sprintf("%s/repos/%s/%s/commits?%s", githubAPIBase, owner, repo, query.Encode())

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getFileContent(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	path, ok := params["path"].(string)
	if !ok {
		return "", fmt.Errorf("path must be a string")
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/contents/%s", githubAPIBase, owner, repo, path)

	if ref, ok := params["ref"].(string); ok && ref != "" {
		endpoint += "?ref=" + url.QueryEscape(ref)
	}

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	// Try to decode base64 content
	var result map[string]interface{}
	if err := httpclient.UnmarshalJSON(respBody, &result); err == nil {
		if content, ok := result["content"].(string); ok {
			if encoding, ok := result["encoding"].(string); ok && encoding == "base64" {
				decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(content, "\n", ""))
				if err == nil {
					result["content"] = string(decoded)
					result["encoding"] = "utf-8"
				}
			}
		}
		return httpclient.PrettyJSONFromInterface(result), nil
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Issues
// =============================================================================

func listIssues(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	query := url.Values{}

	if state, ok := params["state"].(string); ok && state != "" {
		query.Set("state", state)
	} else {
		query.Set("state", "open")
	}

	perPage := 30
	if pp, ok := params["per_page"].(float64); ok {
		perPage = int(pp)
	}
	query.Set("per_page", fmt.Sprintf("%d", perPage))

	page := 1
	if p, ok := params["page"].(float64); ok {
		page = int(p)
	}
	query.Set("page", fmt.Sprintf("%d", page))

	endpoint := fmt.Sprintf("%s/repos/%s/%s/issues?%s", githubAPIBase, owner, repo, query.Encode())

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getIssue(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	issueNumber, ok := params["issue_number"].(float64)
	if !ok {
		return "", fmt.Errorf("issue_number must be a number")
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/issues/%d", githubAPIBase, owner, repo, int(issueNumber))

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func createIssue(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	title, ok := params["title"].(string)
	if !ok {
		return "", fmt.Errorf("title must be a string")
	}

	body := map[string]interface{}{
		"title": title,
	}

	if b, ok := params["body"].(string); ok {
		body["body"] = b
	}

	if labels, ok := params["labels"].([]interface{}); ok {
		body["labels"] = labels
	}

	if assignees, ok := params["assignees"].([]interface{}); ok {
		body["assignees"] = assignees
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/issues", githubAPIBase, owner, repo)

	respBody, err := client.DoJSON("POST", endpoint, headers(), body)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func updateIssue(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	issueNumber, ok := params["issue_number"].(float64)
	if !ok {
		return "", fmt.Errorf("issue_number must be a number")
	}

	body := make(map[string]interface{})

	if title, ok := params["title"].(string); ok {
		body["title"] = title
	}

	if b, ok := params["body"].(string); ok {
		body["body"] = b
	}

	if state, ok := params["state"].(string); ok {
		body["state"] = state
	}

	if labels, ok := params["labels"].([]interface{}); ok {
		body["labels"] = labels
	}

	if assignees, ok := params["assignees"].([]interface{}); ok {
		body["assignees"] = assignees
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/issues/%d", githubAPIBase, owner, repo, int(issueNumber))

	respBody, err := client.DoJSON("PATCH", endpoint, headers(), body)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func addIssueComment(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	issueNumber, ok := params["issue_number"].(float64)
	if !ok {
		return "", fmt.Errorf("issue_number must be a number")
	}

	body, ok := params["body"].(string)
	if !ok {
		return "", fmt.Errorf("body must be a string")
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments", githubAPIBase, owner, repo, int(issueNumber))

	respBody, err := client.DoJSON("POST", endpoint, headers(), map[string]string{"body": body})
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Pull Requests
// =============================================================================

func listPRs(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	query := url.Values{}

	if state, ok := params["state"].(string); ok && state != "" {
		query.Set("state", state)
	} else {
		query.Set("state", "open")
	}

	perPage := 30
	if pp, ok := params["per_page"].(float64); ok {
		perPage = int(pp)
	}
	query.Set("per_page", fmt.Sprintf("%d", perPage))

	page := 1
	if p, ok := params["page"].(float64); ok {
		page = int(p)
	}
	query.Set("page", fmt.Sprintf("%d", page))

	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls?%s", githubAPIBase, owner, repo, query.Encode())

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getPR(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	prNumber, ok := params["pr_number"].(float64)
	if !ok {
		return "", fmt.Errorf("pr_number must be a number")
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", githubAPIBase, owner, repo, int(prNumber))

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func createPR(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	title, ok := params["title"].(string)
	if !ok {
		return "", fmt.Errorf("title must be a string")
	}

	head, ok := params["head"].(string)
	if !ok {
		return "", fmt.Errorf("head must be a string")
	}

	base, ok := params["base"].(string)
	if !ok {
		return "", fmt.Errorf("base must be a string")
	}

	body := map[string]interface{}{
		"title": title,
		"head":  head,
		"base":  base,
	}

	if b, ok := params["body"].(string); ok {
		body["body"] = b
	}

	if draft, ok := params["draft"].(bool); ok {
		body["draft"] = draft
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls", githubAPIBase, owner, repo)

	respBody, err := client.DoJSON("POST", endpoint, headers(), body)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func listPRCommits(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	prNumber, ok := params["pr_number"].(float64)
	if !ok {
		return "", fmt.Errorf("pr_number must be a number")
	}

	perPage := 30
	if pp, ok := params["per_page"].(float64); ok {
		perPage = int(pp)
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/commits?per_page=%d", githubAPIBase, owner, repo, int(prNumber), perPage)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func listPRFiles(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	prNumber, ok := params["pr_number"].(float64)
	if !ok {
		return "", fmt.Errorf("pr_number must be a number")
	}

	perPage := 30
	if pp, ok := params["per_page"].(float64); ok {
		perPage = int(pp)
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/files?per_page=%d", githubAPIBase, owner, repo, int(prNumber), perPage)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func listPRReviews(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	prNumber, ok := params["pr_number"].(float64)
	if !ok {
		return "", fmt.Errorf("pr_number must be a number")
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/reviews", githubAPIBase, owner, repo, int(prNumber))

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Search
// =============================================================================

func searchRepos(params map[string]interface{}) (string, error) {
	query, ok := params["query"].(string)
	if !ok {
		return "", fmt.Errorf("query must be a string")
	}

	q := url.Values{}
	q.Set("q", query)

	if sort, ok := params["sort"].(string); ok && sort != "" {
		q.Set("sort", sort)
	}

	perPage := 30
	if pp, ok := params["per_page"].(float64); ok {
		perPage = int(pp)
	}
	q.Set("per_page", fmt.Sprintf("%d", perPage))

	page := 1
	if p, ok := params["page"].(float64); ok {
		page = int(p)
	}
	q.Set("page", fmt.Sprintf("%d", page))

	endpoint := fmt.Sprintf("%s/search/repositories?%s", githubAPIBase, q.Encode())

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func searchCode(params map[string]interface{}) (string, error) {
	query, ok := params["query"].(string)
	if !ok {
		return "", fmt.Errorf("query must be a string")
	}

	q := url.Values{}
	q.Set("q", query)

	perPage := 30
	if pp, ok := params["per_page"].(float64); ok {
		perPage = int(pp)
	}
	q.Set("per_page", fmt.Sprintf("%d", perPage))

	page := 1
	if p, ok := params["page"].(float64); ok {
		page = int(p)
	}
	q.Set("page", fmt.Sprintf("%d", page))

	endpoint := fmt.Sprintf("%s/search/code?%s", githubAPIBase, q.Encode())

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func searchIssues(params map[string]interface{}) (string, error) {
	query, ok := params["query"].(string)
	if !ok {
		return "", fmt.Errorf("query must be a string")
	}

	q := url.Values{}
	q.Set("q", query)

	if sort, ok := params["sort"].(string); ok && sort != "" {
		q.Set("sort", sort)
	}

	perPage := 30
	if pp, ok := params["per_page"].(float64); ok {
		perPage = int(pp)
	}
	q.Set("per_page", fmt.Sprintf("%d", perPage))

	page := 1
	if p, ok := params["page"].(float64); ok {
		page = int(p)
	}
	q.Set("page", fmt.Sprintf("%d", page))

	endpoint := fmt.Sprintf("%s/search/issues?%s", githubAPIBase, q.Encode())

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Actions
// =============================================================================

func listWorkflows(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	perPage := 30
	if pp, ok := params["per_page"].(float64); ok {
		perPage = int(pp)
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/actions/workflows?per_page=%d", githubAPIBase, owner, repo, perPage)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func listWorkflowRuns(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	query := url.Values{}

	if status, ok := params["status"].(string); ok && status != "" {
		query.Set("status", status)
	}

	perPage := 30
	if pp, ok := params["per_page"].(float64); ok {
		perPage = int(pp)
	}
	query.Set("per_page", fmt.Sprintf("%d", perPage))

	var endpoint string
	if workflowID, ok := params["workflow_id"].(string); ok && workflowID != "" {
		endpoint = fmt.Sprintf("%s/repos/%s/%s/actions/workflows/%s/runs?%s", githubAPIBase, owner, repo, workflowID, query.Encode())
	} else if workflowID, ok := params["workflow_id"].(float64); ok {
		endpoint = fmt.Sprintf("%s/repos/%s/%s/actions/workflows/%d/runs?%s", githubAPIBase, owner, repo, int(workflowID), query.Encode())
	} else {
		endpoint = fmt.Sprintf("%s/repos/%s/%s/actions/runs?%s", githubAPIBase, owner, repo, query.Encode())
	}

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getWorkflowRun(params map[string]interface{}) (string, error) {
	owner, ok := params["owner"].(string)
	if !ok {
		return "", fmt.Errorf("owner must be a string")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", fmt.Errorf("repo must be a string")
	}

	runID, ok := params["run_id"].(float64)
	if !ok {
		return "", fmt.Errorf("run_id must be a number")
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/actions/runs/%d", githubAPIBase, owner, repo, int(runID))

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}
