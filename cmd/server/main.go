package main

import (
	"log"
	"net/http"
	"os"

	"github.com/shibaleo/go-mcp-dev/internal/auth"
	"github.com/shibaleo/go-mcp-dev/internal/mcp"
	"github.com/shibaleo/go-mcp-dev/internal/modules"
	"github.com/shibaleo/go-mcp-dev/internal/modules/supabase"
	"github.com/shibaleo/go-mcp-dev/internal/observability"
)

func main() {
	// Initialize Loki client
	observability.Init()

	// Register modules
	modules.Register(supabase.Module())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := mcp.NewHandler()

	authMiddleware := auth.NewMiddleware(os.Getenv("INTERNAL_SECRET"))

	http.HandleFunc("/health", healthHandler)
	http.Handle("/mcp", authMiddleware(handler))

	log.Printf("Starting MCP server on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
