package supabase

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/shibaleo/go-mcp-dev/internal/httpclient"
	"github.com/shibaleo/go-mcp-dev/internal/modules"
)

const supabaseAPIBase = "https://api.supabase.com/v1"

var client = httpclient.New()

func getAccessToken() string {
	return os.Getenv("SUPABASE_ACCESS_TOKEN")
}

func headers() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + getAccessToken(),
	}
}

// Module returns the Supabase module definition
func Module() modules.ModuleDefinition {
	return modules.ModuleDefinition{
		Name:        "supabase",
		Description: "Supabase Management API - プロジェクト管理、DB操作、マイグレーション、ログ、ストレージ",
		APIVersion:  "v1",
		TestedAt:    "2026-01-10",
		Tools:       tools,
		Handlers:    handlers,
	}
}

var tools = []modules.Tool{
	// Account Tools
	{
		Name:        "list_organizations",
		Description: "List all organizations you have access to.",
		InputSchema: modules.InputSchema{
			Type:       "object",
			Properties: map[string]modules.Property{},
		},
	},
	{
		Name:        "list_projects",
		Description: "List all Supabase projects you have access to. Use this first to get project_ref for other operations.",
		InputSchema: modules.InputSchema{
			Type:       "object",
			Properties: map[string]modules.Property{},
		},
	},
	{
		Name:        "get_project",
		Description: "Get details of a specific project.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference (e.g., 'abcdefghijk'). Get from list_projects.",
				},
			},
			Required: []string{"project_ref"},
		},
	},
	// Database Tools
	{
		Name:        "list_tables",
		Description: "List all tables in the database with their schemas. Returns table names and column counts.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
				"schemas": {
					Type:        "array",
					Description: "Schemas to include (default: ['public'])",
				},
			},
			Required: []string{"project_ref"},
		},
	},
	{
		Name:        "run_query",
		Description: "Execute a SQL query against the database. Supports both read and write operations.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
				"query": {
					Type:        "string",
					Description: "SQL query to execute",
				},
			},
			Required: []string{"project_ref", "query"},
		},
	},
	{
		Name:        "list_migrations",
		Description: "List all database migrations that have been applied.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
			},
			Required: []string{"project_ref"},
		},
	},
	{
		Name:        "apply_migration",
		Description: "Apply a new database migration. Use for DDL operations like CREATE TABLE, ALTER TABLE, etc.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
				"name": {
					Type:        "string",
					Description: "Migration name in snake_case (e.g., add_users_table)",
				},
				"query": {
					Type:        "string",
					Description: "SQL DDL statements to apply",
				},
			},
			Required: []string{"project_ref", "name", "query"},
		},
	},
	// Debugging Tools
	{
		Name:        "get_logs",
		Description: "Get logs for a specific service. Available services: api, postgres, edge-function, auth, storage, realtime.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
				"service": {
					Type:        "string",
					Description: "Service to get logs for: api, postgres, edge-function, auth, storage, realtime",
				},
				"start_time": {
					Type:        "string",
					Description: "ISO timestamp for start of log range (optional)",
				},
				"end_time": {
					Type:        "string",
					Description: "ISO timestamp for end of log range (optional)",
				},
			},
			Required: []string{"project_ref", "service"},
		},
	},
	{
		Name:        "get_security_advisors",
		Description: "Get security recommendations and potential issues for the project.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
			},
			Required: []string{"project_ref"},
		},
	},
	{
		Name:        "get_performance_advisors",
		Description: "Get performance recommendations and potential issues for the project.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
			},
			Required: []string{"project_ref"},
		},
	},
	// Development Tools
	{
		Name:        "get_project_url",
		Description: "Get the base URL for a Supabase project.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
			},
			Required: []string{"project_ref"},
		},
	},
	{
		Name:        "get_api_keys",
		Description: "Get the API keys for the project (anon key and service role key).",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
			},
			Required: []string{"project_ref"},
		},
	},
	{
		Name:        "generate_typescript_types",
		Description: "Generate TypeScript type definitions from the database schema.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
			},
			Required: []string{"project_ref"},
		},
	},
	// Edge Function Tools
	{
		Name:        "list_edge_functions",
		Description: "List all Edge Functions deployed in the project.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
			},
			Required: []string{"project_ref"},
		},
	},
	{
		Name:        "get_edge_function",
		Description: "Get details of a specific Edge Function.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
				"slug": {
					Type:        "string",
					Description: "The slug/name of the Edge Function",
				},
			},
			Required: []string{"project_ref", "slug"},
		},
	},
	// Storage Tools
	{
		Name:        "list_storage_buckets",
		Description: "List all storage buckets in the project.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
			},
			Required: []string{"project_ref"},
		},
	},
	{
		Name:        "get_storage_config",
		Description: "Get storage configuration for the project including file size limits and features.",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"project_ref": {
					Type:        "string",
					Description: "Project reference",
				},
			},
			Required: []string{"project_ref"},
		},
	},
}

var handlers = map[string]modules.ToolHandler{
	"list_organizations":        listOrganizations,
	"list_projects":             listProjects,
	"get_project":               getProject,
	"list_tables":               listTables,
	"run_query":                 runQuery,
	"list_migrations":           listMigrations,
	"apply_migration":           applyMigration,
	"get_logs":                  getLogs,
	"get_security_advisors":     getSecurityAdvisors,
	"get_performance_advisors":  getPerformanceAdvisors,
	"get_project_url":           getProjectURL,
	"get_api_keys":              getAPIKeys,
	"generate_typescript_types": generateTypescriptTypes,
	"list_edge_functions":       listEdgeFunctions,
	"get_edge_function":         getEdgeFunction,
	"list_storage_buckets":      listStorageBuckets,
	"get_storage_config":        getStorageConfig,
}

// =============================================================================
// Account Tools
// =============================================================================

func listOrganizations(params map[string]interface{}) (string, error) {
	endpoint := supabaseAPIBase + "/organizations"

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func listProjects(params map[string]interface{}) (string, error) {
	endpoint := supabaseAPIBase + "/projects"

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getProject(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	endpoint := fmt.Sprintf("%s/projects/%s", supabaseAPIBase, projectRef)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Database Tools
// =============================================================================

func listTables(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	// Default schemas
	schemas := []string{"public"}
	if s, ok := params["schemas"].([]interface{}); ok && len(s) > 0 {
		schemas = make([]string, 0, len(s))
		for _, schema := range s {
			if str, ok := schema.(string); ok {
				schemas = append(schemas, str)
			}
		}
	}

	// Build query to list tables
	schemaList := make([]string, len(schemas))
	for i, s := range schemas {
		schemaList[i] = fmt.Sprintf("'%s'", s)
	}

	query := fmt.Sprintf(`
		SELECT
			schemaname as schema,
			tablename as name,
			(SELECT count(*)::int FROM information_schema.columns
			 WHERE table_schema = t.schemaname AND table_name = t.tablename) as column_count
		FROM pg_tables t
		WHERE schemaname = ANY(ARRAY[%s])
		ORDER BY schemaname, tablename
	`, strings.Join(schemaList, ","))

	return executeQuery(projectRef, query)
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

	return executeQuery(projectRef, query)
}

func executeQuery(projectRef, query string) (string, error) {
	endpoint := fmt.Sprintf("%s/projects/%s/database/query", supabaseAPIBase, projectRef)

	payload := map[string]string{"query": query}

	respBody, err := client.DoJSON("POST", endpoint, headers(), payload)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func listMigrations(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	// Query to list migrations from supabase_migrations schema
	query := `
		SELECT version, name, executed_at
		FROM supabase_migrations.schema_migrations
		ORDER BY version DESC
	`

	return executeQuery(projectRef, query)
}

func applyMigration(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	name, ok := params["name"].(string)
	if !ok {
		return "", fmt.Errorf("name must be a string")
	}

	query, ok := params["query"].(string)
	if !ok {
		return "", fmt.Errorf("query must be a string")
	}

	// Execute the migration
	_, err := executeQuery(projectRef, query)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`{"success": true, "migration": "%s"}`, name), nil
}

// =============================================================================
// Debugging Tools
// =============================================================================

func getLogs(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	service, ok := params["service"].(string)
	if !ok {
		return "", fmt.Errorf("service must be a string")
	}

	// Map service names to log collection names
	serviceMap := map[string]string{
		"api":           "api_logs",
		"postgres":      "postgres_logs",
		"edge-function": "function_logs",
		"auth":          "auth_logs",
		"storage":       "storage_logs",
		"realtime":      "realtime_logs",
	}

	collection, exists := serviceMap[service]
	if !exists {
		return "", fmt.Errorf("invalid service: %s. Valid services: api, postgres, edge-function, auth, storage, realtime", service)
	}

	query := url.Values{}
	query.Set("collection", collection)

	if startTime, ok := params["start_time"].(string); ok && startTime != "" {
		query.Set("start", startTime)
	}

	if endTime, ok := params["end_time"].(string); ok && endTime != "" {
		query.Set("end", endTime)
	}

	endpoint := fmt.Sprintf("%s/projects/%s/analytics/endpoints/logs.all?%s", supabaseAPIBase, projectRef, query.Encode())

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getSecurityAdvisors(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/advisors/security", supabaseAPIBase, projectRef)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getPerformanceAdvisors(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/advisors/performance", supabaseAPIBase, projectRef)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Development Tools
// =============================================================================

func getProjectURL(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	projectURL := fmt.Sprintf("https://%s.supabase.co", projectRef)

	return fmt.Sprintf(`{"url": "%s"}`, projectURL), nil
}

func getAPIKeys(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/api-keys", supabaseAPIBase, projectRef)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func generateTypescriptTypes(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/types/typescript", supabaseAPIBase, projectRef)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Edge Function Tools
// =============================================================================

func listEdgeFunctions(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/functions", supabaseAPIBase, projectRef)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getEdgeFunction(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	slug, ok := params["slug"].(string)
	if !ok {
		return "", fmt.Errorf("slug must be a string")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/functions/%s", supabaseAPIBase, projectRef, slug)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Storage Tools
// =============================================================================

func listStorageBuckets(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/storage/buckets", supabaseAPIBase, projectRef)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getStorageConfig(params map[string]interface{}) (string, error) {
	projectRef, ok := params["project_ref"].(string)
	if !ok {
		return "", fmt.Errorf("project_ref must be a string")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/config/storage", supabaseAPIBase, projectRef)

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}
