package app

// These are set at buildtime - see makefile
var (
	Name      = "unknown"
	Product   = "unknown"
	Branch    = "unknown"
	BuildDate = "unknown"
	Commit    = "unknown"
	Version   = "unknown"
	Author    = "unknown"
)

// SourceDetails is useful for debugging
type SourceDetails struct {
	Name      string `json:"name,omitempty"`
	Product   string `json:"product,omitempty"`
	Branch    string `json:"branch,omitempty"`
	BuildDate string `json:"build_date,omitempty"`
	Commit    string `json:"commit,omitempty"`
	Version   string `json:"version,omitempty"`
	Author    string `json:"author,omitempty"`
}

// New returns a SourceDetails
func New() *SourceDetails {
	return &SourceDetails{
		Name:      Name,
		Product:   Product,
		Branch:    Branch,
		BuildDate: BuildDate,
		Commit:    Commit,
		Version:   Version,
		Author:    Author,
	}
}
