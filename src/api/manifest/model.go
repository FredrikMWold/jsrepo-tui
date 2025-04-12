package manifest

type Block struct {
	Name              string            `json:"name"`
	Directory         string            `json:"directory"`
	Category          string            `json:"category"`
	Tests             bool              `json:"tests"`
	Subdirectory      bool              `json:"subdirectory"`
	List              bool              `json:"list"`
	Files             []string          `json:"files"`
	LocalDependencies []string          `json:"localDependencies"`
	Dependencies      []string          `json:"dependencies"`
	DevDependencies   []string          `json:"devDependencies"`
	Imports           map[string]string `json:"_imports_"`
}

type Category struct {
	Name   string  `json:"name"`
	Blocks []Block `json:"blocks"`
}

type ManifestResponse struct {
	Categories []Category `json:"categories"`
}
