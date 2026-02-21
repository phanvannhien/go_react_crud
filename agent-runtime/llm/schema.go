package llm

type PlanResult struct {
	Summary       string   `json:"summary"`
	FilesToModify []string `json:"filesToModify"`
	RiskScore     float64  `json:"riskScore"`
}

type FileChange struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type ImplementationResult struct {
	Files []FileChange `json:"files"`
}
