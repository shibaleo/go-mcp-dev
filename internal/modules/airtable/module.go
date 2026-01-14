package airtable

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/shibaleo/go-mcp-dev/internal/httpclient"
	"github.com/shibaleo/go-mcp-dev/internal/modules"
)

const (
	airtableAPIBase = "https://api.airtable.com/v0"
	airtableVersion = "v0"
)

var client = httpclient.New()

func getToken() string {
	return os.Getenv("AIRTABLE_API_KEY")
}

func headers() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + getToken(),
		"Content-Type":  "application/json",
	}
}

// Module returns the Airtable module definition
func Module() modules.ModuleDefinition {
	return modules.ModuleDefinition{
		Name:        "airtable",
		Description: "Airtable API - Bases, Tables, Records operations",
		APIVersion:  airtableVersion,
		TestedAt:    "2026-01-14",
		Tools:       tools,
		Handlers:    handlers,
	}
}

var tools = []modules.Tool{
	// Base Operations
	{
		Name:        "list_bases",
		Description: "List all accessible Airtable bases with their names, IDs, and permission levels",
		InputSchema: modules.InputSchema{
			Type:       "object",
			Properties: map[string]modules.Property{},
		},
	},
	// Schema Operations
	{
		Name:        "describe",
		Description: "Describe Airtable base or table schema. Use detailLevel to optimize context: tableIdentifiersOnly (minimal), identifiersOnly (IDs and names), full (complete details with field types)",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"base_id": {
					Type:        "string",
					Description: "Base ID (starts with 'app')",
				},
				"scope": {
					Type:        "string",
					Description: "Scope of description: 'base' for all tables, 'table' for a specific table",
				},
				"table": {
					Type:        "string",
					Description: "Table name or ID (required when scope='table')",
				},
				"detail_level": {
					Type:        "string",
					Description: "Detail level: tableIdentifiersOnly, identifiersOnly, or full (default: full)",
				},
				"include_views": {
					Type:        "boolean",
					Description: "Include view information (default: false)",
				},
			},
			Required: []string{"base_id"},
		},
	},
	// Record Operations
	{
		Name:        "query",
		Description: "Query Airtable records with filtering, sorting, and pagination",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"base_id": {
					Type:        "string",
					Description: "Base ID (starts with 'app')",
				},
				"table": {
					Type:        "string",
					Description: "Table name or ID",
				},
				"fields": {
					Type:        "array",
					Description: "Array of field names to return",
				},
				"filter_by_formula": {
					Type:        "string",
					Description: "Airtable formula to filter records",
				},
				"view": {
					Type:        "string",
					Description: "View name or ID to use",
				},
				"sort": {
					Type:        "array",
					Description: "Sort configuration: [{field: string, direction: 'asc'|'desc'}]",
				},
				"page_size": {
					Type:        "number",
					Description: "Number of records per page (max 100)",
				},
				"max_records": {
					Type:        "number",
					Description: "Maximum number of records to return",
				},
				"offset": {
					Type:        "string",
					Description: "Pagination offset from previous response",
				},
			},
			Required: []string{"base_id", "table"},
		},
	},
	{
		Name:        "get_record",
		Description: "Retrieve a single record by ID",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"base_id": {
					Type:        "string",
					Description: "Base ID (starts with 'app')",
				},
				"table": {
					Type:        "string",
					Description: "Table name or ID",
				},
				"record_id": {
					Type:        "string",
					Description: "Record ID (starts with 'rec')",
				},
			},
			Required: []string{"base_id", "table", "record_id"},
		},
	},
	{
		Name:        "create",
		Description: "Create new records in a table. Supports batch creation (up to 10 records per request)",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"base_id": {
					Type:        "string",
					Description: "Base ID (starts with 'app')",
				},
				"table": {
					Type:        "string",
					Description: "Table name or ID",
				},
				"records": {
					Type:        "array",
					Description: "Array of records to create. Each record is {fields: {fieldName: value}}",
				},
				"typecast": {
					Type:        "boolean",
					Description: "Automatically typecast field values (default: false)",
				},
			},
			Required: []string{"base_id", "table", "records"},
		},
	},
	{
		Name:        "update",
		Description: "Update existing records. Supports batch update (up to 10 records per request). Uses PATCH (partial update)",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"base_id": {
					Type:        "string",
					Description: "Base ID (starts with 'app')",
				},
				"table": {
					Type:        "string",
					Description: "Table name or ID",
				},
				"records": {
					Type:        "array",
					Description: "Array of records to update. Each record is {id: recordId, fields: {fieldName: value}}",
				},
				"typecast": {
					Type:        "boolean",
					Description: "Automatically typecast field values (default: false)",
				},
			},
			Required: []string{"base_id", "table", "records"},
		},
	},
	{
		Name:        "delete",
		Description: "Delete records from a table. Supports batch deletion (up to 10 records per request)",
		InputSchema: modules.InputSchema{
			Type: "object",
			Properties: map[string]modules.Property{
				"base_id": {
					Type:        "string",
					Description: "Base ID (starts with 'app')",
				},
				"table": {
					Type:        "string",
					Description: "Table name or ID",
				},
				"record_ids": {
					Type:        "array",
					Description: "Array of record IDs to delete",
				},
			},
			Required: []string{"base_id", "table", "record_ids"},
		},
	},
}

var handlers = map[string]modules.ToolHandler{
	"list_bases": listBases,
	"describe":   describe,
	"query":      query,
	"get_record": getRecord,
	"create":     create,
	"update":     update,
	"delete":     deleteRecords,
}

// =============================================================================
// Base Operations
// =============================================================================

func listBases(params map[string]interface{}) (string, error) {
	endpoint := "https://api.airtable.com/v0/meta/bases"

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

// =============================================================================
// Schema Operations
// =============================================================================

func describe(params map[string]interface{}) (string, error) {
	baseID, ok := params["base_id"].(string)
	if !ok || baseID == "" {
		return "", fmt.Errorf("base_id is required")
	}

	scope, _ := params["scope"].(string)
	if scope == "" {
		scope = "base"
	}

	// Get tables (this endpoint returns table schema)
	tablesEndpoint := fmt.Sprintf("https://api.airtable.com/v0/meta/bases/%s/tables", url.PathEscape(baseID))
	tablesInfoBytes, err := client.DoJSON("GET", tablesEndpoint, headers(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to get tables: %w", err)
	}

	var tablesData map[string]interface{}
	if err := json.Unmarshal(tablesInfoBytes, &tablesData); err != nil {
		return "", fmt.Errorf("failed to parse tables: %w", err)
	}

	// Build response
	result := map[string]interface{}{
		"base_id": baseID,
	}
	tables, _ := tablesData["tables"].([]interface{})

	detailLevel, _ := params["detail_level"].(string)
	if detailLevel == "" {
		detailLevel = "full"
	}

	includeViews := false
	if iv, ok := params["include_views"].(bool); ok {
		includeViews = iv
	}

	// Filter to specific table if scope is "table"
	if scope == "table" {
		tableName, _ := params["table"].(string)
		if tableName == "" {
			return "", fmt.Errorf("table is required when scope is 'table'")
		}

		var foundTable interface{}
		for _, t := range tables {
			tbl, _ := t.(map[string]interface{})
			if tbl["name"] == tableName || tbl["id"] == tableName {
				foundTable = t
				break
			}
		}

		if foundTable == nil {
			return "", fmt.Errorf("table '%s' not found", tableName)
		}

		tables = []interface{}{foundTable}
	}

	// Apply detail level filtering
	filteredTables := filterTablesDetail(tables, detailLevel, includeViews)
	result["tables"] = filteredTables

	return httpclient.PrettyJSONFromInterface(result), nil
}

func filterTablesDetail(tables []interface{}, detailLevel string, includeViews bool) []interface{} {
	result := make([]interface{}, 0, len(tables))

	for _, t := range tables {
		tbl, ok := t.(map[string]interface{})
		if !ok {
			continue
		}

		filtered := make(map[string]interface{})
		filtered["id"] = tbl["id"]
		filtered["name"] = tbl["name"]

		switch detailLevel {
		case "tableIdentifiersOnly":
			// Only id and name

		case "identifiersOnly":
			// Add primaryFieldId
			if pf, ok := tbl["primaryFieldId"]; ok {
				filtered["primaryFieldId"] = pf
			}
			// Add field identifiers only
			if fields, ok := tbl["fields"].([]interface{}); ok {
				filteredFields := make([]map[string]interface{}, 0, len(fields))
				for _, f := range fields {
					field, _ := f.(map[string]interface{})
					filteredFields = append(filteredFields, map[string]interface{}{
						"id":   field["id"],
						"name": field["name"],
					})
				}
				filtered["fields"] = filteredFields
			}
			// Add view identifiers if requested
			if includeViews {
				if views, ok := tbl["views"].([]interface{}); ok {
					filteredViews := make([]map[string]interface{}, 0, len(views))
					for _, v := range views {
						view, _ := v.(map[string]interface{})
						filteredViews = append(filteredViews, map[string]interface{}{
							"id":   view["id"],
							"name": view["name"],
						})
					}
					filtered["views"] = filteredViews
				}
			}

		default: // "full"
			// Include everything
			if pf, ok := tbl["primaryFieldId"]; ok {
				filtered["primaryFieldId"] = pf
			}
			if fields, ok := tbl["fields"]; ok {
				filtered["fields"] = fields
			}
			if includeViews {
				if views, ok := tbl["views"]; ok {
					filtered["views"] = views
				}
			}
		}

		result = append(result, filtered)
	}

	return result
}

// =============================================================================
// Record Operations
// =============================================================================

func query(params map[string]interface{}) (string, error) {
	baseID, ok := params["base_id"].(string)
	if !ok || baseID == "" {
		return "", fmt.Errorf("base_id is required")
	}

	table, ok := params["table"].(string)
	if !ok || table == "" {
		return "", fmt.Errorf("table is required")
	}

	// Build query parameters
	queryParams := url.Values{}

	if fields, ok := params["fields"].([]interface{}); ok {
		for _, f := range fields {
			if fieldName, ok := f.(string); ok {
				queryParams.Add("fields[]", fieldName)
			}
		}
	}

	if formula, ok := params["filter_by_formula"].(string); ok && formula != "" {
		queryParams.Set("filterByFormula", formula)
	}

	if view, ok := params["view"].(string); ok && view != "" {
		queryParams.Set("view", view)
	}

	if sorts, ok := params["sort"].([]interface{}); ok {
		for i, s := range sorts {
			sort, _ := s.(map[string]interface{})
			if field, ok := sort["field"].(string); ok {
				queryParams.Set(fmt.Sprintf("sort[%d][field]", i), field)
				direction := "asc"
				if dir, ok := sort["direction"].(string); ok {
					direction = dir
				}
				queryParams.Set(fmt.Sprintf("sort[%d][direction]", i), direction)
			}
		}
	}

	if pageSize, ok := params["page_size"].(float64); ok {
		queryParams.Set("pageSize", fmt.Sprintf("%d", int(pageSize)))
	}

	if maxRecords, ok := params["max_records"].(float64); ok {
		queryParams.Set("maxRecords", fmt.Sprintf("%d", int(maxRecords)))
	}

	if offset, ok := params["offset"].(string); ok && offset != "" {
		queryParams.Set("offset", offset)
	}

	endpoint := fmt.Sprintf("%s/%s/%s", airtableAPIBase, url.PathEscape(baseID), url.PathEscape(table))
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func getRecord(params map[string]interface{}) (string, error) {
	baseID, ok := params["base_id"].(string)
	if !ok || baseID == "" {
		return "", fmt.Errorf("base_id is required")
	}

	table, ok := params["table"].(string)
	if !ok || table == "" {
		return "", fmt.Errorf("table is required")
	}

	recordID, ok := params["record_id"].(string)
	if !ok || recordID == "" {
		return "", fmt.Errorf("record_id is required")
	}

	endpoint := fmt.Sprintf("%s/%s/%s/%s", airtableAPIBase, url.PathEscape(baseID), url.PathEscape(table), url.PathEscape(recordID))

	respBody, err := client.DoJSON("GET", endpoint, headers(), nil)
	if err != nil {
		return "", err
	}

	return httpclient.PrettyJSON(respBody), nil
}

func create(params map[string]interface{}) (string, error) {
	baseID, ok := params["base_id"].(string)
	if !ok || baseID == "" {
		return "", fmt.Errorf("base_id is required")
	}

	table, ok := params["table"].(string)
	if !ok || table == "" {
		return "", fmt.Errorf("table is required")
	}

	records, ok := params["records"].([]interface{})
	if !ok || len(records) == 0 {
		return "", fmt.Errorf("records is required and must not be empty")
	}

	typecast := false
	if tc, ok := params["typecast"].(bool); ok {
		typecast = tc
	}

	// Process records in chunks of 10
	allCreated := make([]interface{}, 0)

	for i := 0; i < len(records); i += 10 {
		end := i + 10
		if end > len(records) {
			end = len(records)
		}
		chunk := records[i:end]

		body := map[string]interface{}{
			"records":  chunk,
			"typecast": typecast,
		}

		endpoint := fmt.Sprintf("%s/%s/%s", airtableAPIBase, url.PathEscape(baseID), url.PathEscape(table))

		respBody, err := client.DoJSON("POST", endpoint, headers(), body)
		if err != nil {
			return "", fmt.Errorf("failed to create records (batch %d): %w", i/10+1, err)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return "", fmt.Errorf("failed to parse response: %w", err)
		}
		if created, ok := resp["records"].([]interface{}); ok {
			allCreated = append(allCreated, created...)
		}
	}

	result := map[string]interface{}{
		"records": allCreated,
		"summary": map[string]interface{}{
			"created": len(allCreated),
		},
	}

	return httpclient.PrettyJSONFromInterface(result), nil
}

func update(params map[string]interface{}) (string, error) {
	baseID, ok := params["base_id"].(string)
	if !ok || baseID == "" {
		return "", fmt.Errorf("base_id is required")
	}

	table, ok := params["table"].(string)
	if !ok || table == "" {
		return "", fmt.Errorf("table is required")
	}

	records, ok := params["records"].([]interface{})
	if !ok || len(records) == 0 {
		return "", fmt.Errorf("records is required and must not be empty")
	}

	typecast := false
	if tc, ok := params["typecast"].(bool); ok {
		typecast = tc
	}

	// Process records in chunks of 10
	allUpdated := make([]interface{}, 0)

	for i := 0; i < len(records); i += 10 {
		end := i + 10
		if end > len(records) {
			end = len(records)
		}
		chunk := records[i:end]

		body := map[string]interface{}{
			"records":  chunk,
			"typecast": typecast,
		}

		endpoint := fmt.Sprintf("%s/%s/%s", airtableAPIBase, url.PathEscape(baseID), url.PathEscape(table))

		respBody, err := client.DoJSON("PATCH", endpoint, headers(), body)
		if err != nil {
			return "", fmt.Errorf("failed to update records (batch %d): %w", i/10+1, err)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return "", fmt.Errorf("failed to parse response: %w", err)
		}
		if updated, ok := resp["records"].([]interface{}); ok {
			allUpdated = append(allUpdated, updated...)
		}
	}

	result := map[string]interface{}{
		"records": allUpdated,
		"summary": map[string]interface{}{
			"updated": len(allUpdated),
		},
	}

	return httpclient.PrettyJSONFromInterface(result), nil
}

func deleteRecords(params map[string]interface{}) (string, error) {
	baseID, ok := params["base_id"].(string)
	if !ok || baseID == "" {
		return "", fmt.Errorf("base_id is required")
	}

	table, ok := params["table"].(string)
	if !ok || table == "" {
		return "", fmt.Errorf("table is required")
	}

	recordIDs, ok := params["record_ids"].([]interface{})
	if !ok || len(recordIDs) == 0 {
		return "", fmt.Errorf("record_ids is required and must not be empty")
	}

	// Process records in chunks of 10
	allDeleted := make([]interface{}, 0)

	for i := 0; i < len(recordIDs); i += 10 {
		end := i + 10
		if end > len(recordIDs) {
			end = len(recordIDs)
		}
		chunk := recordIDs[i:end]

		// Build query parameters for DELETE
		queryParams := url.Values{}
		for _, id := range chunk {
			if recordID, ok := id.(string); ok {
				queryParams.Add("records[]", recordID)
			}
		}

		endpoint := fmt.Sprintf("%s/%s/%s?%s", airtableAPIBase, url.PathEscape(baseID), url.PathEscape(table), queryParams.Encode())

		respBody, err := client.DoJSON("DELETE", endpoint, headers(), nil)
		if err != nil {
			return "", fmt.Errorf("failed to delete records (batch %d): %w", i/10+1, err)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return "", fmt.Errorf("failed to parse response: %w", err)
		}
		if deleted, ok := resp["records"].([]interface{}); ok {
			allDeleted = append(allDeleted, deleted...)
		}
	}

	result := map[string]interface{}{
		"records": allDeleted,
		"summary": map[string]interface{}{
			"deleted": len(allDeleted),
		},
	}

	return httpclient.PrettyJSONFromInterface(result), nil
}

// Ensure json package is used
var _ = json.Marshal


