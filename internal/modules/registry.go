package modules

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shibaleo/go-mcp-dev/internal/observability"
)

// Registry holds all module definitions
var Registry = make(map[string]ModuleDefinition)

// Register adds a module to the registry
func Register(module ModuleDefinition) {
	Registry[module.Name] = module
}

// MetaTools returns the two meta tools for lazy loading
func MetaTools() []Tool {
	return []Tool{
		{
			Name: "get_module_schema",
			Description: `モジュールのツール定義を取得。重要: 各モジュールにつき1セッション1回のみ呼び出すこと。スキーマは会話履歴にキャッシュされるため、同一モジュールへの2回目以降の呼び出しはcall_module_toolを直接使用すること。`,
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"module": {
						Type:        "string",
						Description: "モジュール名(google_calendar, microsoft_todo, notion, rag, supabase, jira, confluence, github)",
					},
				},
				Required: []string{"module"},
			},
		},
		{
			Name: "call_module_tool",
			Description: `モジュールのツールを呼び出す。

【利用可能モジュール】
- google_calendar: 予定の取得・作成
- microsoft_todo: タスク管理
- notion: ページ・データベース操作
- rag: ドキュメント検索(セマンティック検索、キーワード検索)
- supabase: DB操作、マイグレーション、ログ、ストレージ
- jira: Issue/Project操作（検索、作成、更新、コメント）
- confluence: Wiki操作（スペース、ページ、検索、ラベル）
- github: リポジトリ、Issue、PR、Actions、検索

【使い方】
1. get_module_schema(module) でツール一覧とパラメータを確認
2. call_module_tool(module, tool_name, params) で実行`,
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"module": {
						Type:        "string",
						Description: "モジュール名",
					},
					"tool_name": {
						Type:        "string",
						Description: "ツール名",
					},
					"params": {
						Type:        "object",
						Description: "ツールパラメータ",
					},
				},
				Required: []string{"module", "tool_name"},
			},
		},
	}
}

// GetModuleSchema returns the schema for a module
func GetModuleSchema(moduleName string) (*ToolCallResult, error) {
	module, ok := Registry[moduleName]
	if !ok {
		return &ToolCallResult{
			Content: []ContentBlock{{Type: "text", Text: fmt.Sprintf("Unknown module: %s", moduleName)}},
			IsError: true,
		}, nil
	}

	schema := struct {
		Module      string `json:"module"`
		Description string `json:"description"`
		Tools       []Tool `json:"tools"`
	}{
		Module:      moduleName,
		Description: module.Description,
		Tools:       module.Tools,
	}

	jsonBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return nil, err
	}

	return &ToolCallResult{
		Content: []ContentBlock{{Type: "text", Text: string(jsonBytes)}},
	}, nil
}

// CallModuleTool executes a tool in a module
func CallModuleTool(moduleName, toolName string, params map[string]interface{}) (*ToolCallResult, error) {
	start := time.Now()

	module, ok := Registry[moduleName]
	if !ok {
		return &ToolCallResult{
			Content: []ContentBlock{{Type: "text", Text: fmt.Sprintf("Unknown module: %s", moduleName)}},
			IsError: true,
		}, nil
	}

	handler, ok := module.Handlers[toolName]
	if !ok {
		return &ToolCallResult{
			Content: []ContentBlock{{Type: "text", Text: fmt.Sprintf("Unknown tool: %s in module %s", toolName, moduleName)}},
			IsError: true,
		}, nil
	}

	result, err := handler(params)
	durationMs := time.Since(start).Milliseconds()

	if err != nil {
		observability.LogToolCall(moduleName, toolName, durationMs, "error", err.Error())
		return &ToolCallResult{
			Content: []ContentBlock{{Type: "text", Text: err.Error()}},
			IsError: true,
		}, nil
	}

	observability.LogToolCall(moduleName, toolName, durationMs, "success", "")
	return &ToolCallResult{
		Content: []ContentBlock{{Type: "text", Text: result}},
	}, nil
}
