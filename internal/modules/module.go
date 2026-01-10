package modules

// Info contains metadata about a module
type Info struct {
	Name       string // Module name (e.g., "supabase", "github")
	APIVersion string // API version this module is tested against
	TestedAt   string // Date when last verified (YYYY-MM-DD)
}

// Registry holds all module information
var Registry = []Info{
	{
		Name:       "supabase",
		APIVersion: "v1",
		TestedAt:   "2026-01-10",
	},
}
