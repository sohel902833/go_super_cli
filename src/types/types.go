package types

type Module struct {
	ModuleName      string `json:"moduleName"`
	ModelProperties string `json:"modelProperties"`
}

type Field struct {
	Name     string
	Type     string
	Required bool
}

type FileInstruction struct {
	FilePath    string `json:"filePath"`
	Content     string `json:"content"`
	Description string `json:"description,omitempty"`
}

type UpdateInstruction struct {
	FilePath          string `json:"filePath"`
	Placeholder       string `json:"placeholder"`
	Content           string `json:"content"`
	Position          string `json:"position"` // "top", "bottom", or empty (default top)
	CreateIfNotExists bool   `json:"createIfNotExists"`
	Description       string `json:"description,omitempty"`
}

type ProjectConfig struct {
	Name               string              `json:"name"`
	Version            string              `json:"version"`
	FileInstructions   []FileInstruction   `json:"fileInstructions"`
	UpdateInstructions []UpdateInstruction `json:"updateInstructions"`
}