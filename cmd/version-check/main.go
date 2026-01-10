// Package main provides a CLI tool to verify API versions for all modules.
//
// Each module declares the external API's official version string (NOT semver):
//   - Supabase: "v1" (URL path)
//   - Notion: "2022-06-28" (Notion-Version header)
//   - GitHub: "2022-11-28" (X-GitHub-Api-Version header)
//   - Jira: "3" (URL path /rest/api/3)
//   - Confluence: "v2" (URL path /wiki/api/v2)
//
// This tool sends real requests to each API to verify compatibility.
// Run in CI (main push) to detect breaking changes before deployment.
//
// Exit codes:
//   - 0: All checks passed or skipped
//   - 1: API error (network, auth, etc.)
//   - 2: Version mismatch detected
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/shibaleo/go-mcp-dev/internal/modules"
	"github.com/shibaleo/go-mcp-dev/internal/modules/confluence"
	"github.com/shibaleo/go-mcp-dev/internal/modules/github"
	"github.com/shibaleo/go-mcp-dev/internal/modules/jira"
	"github.com/shibaleo/go-mcp-dev/internal/modules/notion"
	"github.com/shibaleo/go-mcp-dev/internal/modules/supabase"
)

type VersionChecker struct {
	client *http.Client
}

type CheckResult struct {
	Module     string `json:"module"`
	Expected   string `json:"expected"`
	Actual     string `json:"actual"`
	TestedAt   string `json:"tested_at"`
	Status     string `json:"status"` // "ok", "mismatch", "error", "skip"
	Message    string `json:"message,omitempty"`
}

func NewVersionChecker() *VersionChecker {
	return &VersionChecker{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (vc *VersionChecker) CheckAll() []CheckResult {
	results := []CheckResult{}

	// Register modules
	allModules := []modules.ModuleDefinition{
		supabase.Module(),
		notion.Module(),
		github.Module(),
		jira.Module(),
		confluence.Module(),
	}

	for _, mod := range allModules {
		result := vc.checkModule(mod)
		results = append(results, result)
	}

	return results
}

func (vc *VersionChecker) checkModule(mod modules.ModuleDefinition) CheckResult {
	switch mod.Name {
	case "github":
		return vc.checkGitHub(mod)
	case "notion":
		return vc.checkNotion(mod)
	case "supabase":
		return vc.checkSupabase(mod)
	case "jira":
		return vc.checkJira(mod)
	case "confluence":
		return vc.checkConfluence(mod)
	default:
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "skip",
			Message:  "No version checker implemented",
		}
	}
}

func (vc *VersionChecker) checkGitHub(mod modules.ModuleDefinition) CheckResult {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "skip",
			Message:  "GITHUB_TOKEN not set",
		}
	}

	req, _ := http.NewRequest("GET", "https://api.github.com/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", mod.APIVersion)

	resp, err := vc.client.Do(req)
	if err != nil {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "error",
			Message:  err.Error(),
		}
	}
	defer resp.Body.Close()

	// GitHub returns the API version in response headers
	// If we get a 200, our requested version is supported
	if resp.StatusCode == 200 {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			Actual:   mod.APIVersion, // The version we requested works
			TestedAt: mod.TestedAt,
			Status:   "ok",
		}
	}

	return CheckResult{
		Module:   mod.Name,
		Expected: mod.APIVersion,
		TestedAt: mod.TestedAt,
		Status:   "error",
		Message:  fmt.Sprintf("API returned status %d", resp.StatusCode),
	}
}

func (vc *VersionChecker) checkNotion(mod modules.ModuleDefinition) CheckResult {
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "skip",
			Message:  "NOTION_TOKEN not set",
		}
	}

	req, _ := http.NewRequest("GET", "https://api.notion.com/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Notion-Version", mod.APIVersion)

	resp, err := vc.client.Do(req)
	if err != nil {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "error",
			Message:  err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			Actual:   mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "ok",
		}
	}

	// Check for deprecation warning
	if resp.StatusCode == 400 {
		body, _ := io.ReadAll(resp.Body)
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "mismatch",
			Message:  string(body),
		}
	}

	return CheckResult{
		Module:   mod.Name,
		Expected: mod.APIVersion,
		TestedAt: mod.TestedAt,
		Status:   "error",
		Message:  fmt.Sprintf("API returned status %d", resp.StatusCode),
	}
}

func (vc *VersionChecker) checkSupabase(mod modules.ModuleDefinition) CheckResult {
	token := os.Getenv("SUPABASE_ACCESS_TOKEN")
	if token == "" {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "skip",
			Message:  "SUPABASE_ACCESS_TOKEN not set",
		}
	}

	// Supabase doesn't have explicit version headers, just check API is accessible
	req, _ := http.NewRequest("GET", "https://api.supabase.com/v1/projects", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := vc.client.Do(req)
	if err != nil {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "error",
			Message:  err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			Actual:   mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "ok",
		}
	}

	return CheckResult{
		Module:   mod.Name,
		Expected: mod.APIVersion,
		TestedAt: mod.TestedAt,
		Status:   "error",
		Message:  fmt.Sprintf("API returned status %d", resp.StatusCode),
	}
}

func (vc *VersionChecker) checkJira(mod modules.ModuleDefinition) CheckResult {
	domain := os.Getenv("JIRA_DOMAIN")
	email := os.Getenv("JIRA_EMAIL")
	token := os.Getenv("JIRA_API_TOKEN")

	if domain == "" || email == "" || token == "" {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "skip",
			Message:  "JIRA credentials not set",
		}
	}

	// Check the API version endpoint
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/rest/api/3/serverInfo", domain), nil)
	req.SetBasicAuth(email, token)

	resp, err := vc.client.Do(req)
	if err != nil {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "error",
			Message:  err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var info map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&info)
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			Actual:   mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "ok",
		}
	}

	return CheckResult{
		Module:   mod.Name,
		Expected: mod.APIVersion,
		TestedAt: mod.TestedAt,
		Status:   "error",
		Message:  fmt.Sprintf("API returned status %d", resp.StatusCode),
	}
}

func (vc *VersionChecker) checkConfluence(mod modules.ModuleDefinition) CheckResult {
	domain := os.Getenv("CONFLUENCE_DOMAIN")
	email := os.Getenv("CONFLUENCE_EMAIL")
	token := os.Getenv("CONFLUENCE_API_TOKEN")

	if domain == "" || email == "" || token == "" {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "skip",
			Message:  "CONFLUENCE credentials not set",
		}
	}

	// Check v2 API is accessible
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/wiki/api/v2/spaces?limit=1", domain), nil)
	req.SetBasicAuth(email, token)

	resp, err := vc.client.Do(req)
	if err != nil {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "error",
			Message:  err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return CheckResult{
			Module:   mod.Name,
			Expected: mod.APIVersion,
			Actual:   mod.APIVersion,
			TestedAt: mod.TestedAt,
			Status:   "ok",
		}
	}

	return CheckResult{
		Module:   mod.Name,
		Expected: mod.APIVersion,
		TestedAt: mod.TestedAt,
		Status:   "error",
		Message:  fmt.Sprintf("API returned status %d", resp.StatusCode),
	}
}

func main() {
	checker := NewVersionChecker()
	results := checker.CheckAll()

	// Output as JSON
	output, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(output))

	// Check for failures
	hasError := false
	hasMismatch := false

	for _, r := range results {
		switch r.Status {
		case "error":
			hasError = true
			fmt.Fprintf(os.Stderr, "ERROR: %s - %s\n", r.Module, r.Message)
		case "mismatch":
			hasMismatch = true
			fmt.Fprintf(os.Stderr, "MISMATCH: %s - expected %s, got %s\n", r.Module, r.Expected, r.Actual)
		case "ok":
			fmt.Fprintf(os.Stderr, "OK: %s (%s, tested %s)\n", r.Module, r.Expected, r.TestedAt)
		case "skip":
			fmt.Fprintf(os.Stderr, "SKIP: %s - %s\n", r.Module, r.Message)
		}
	}

	if hasMismatch {
		os.Exit(2)
	}
	if hasError {
		os.Exit(1)
	}
}
