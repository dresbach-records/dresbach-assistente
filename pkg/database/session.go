package database

// PreAnalysisData armazena os dados coletados na fase de pré-análise.
type PreAnalysisData struct {
	RepoURL            string `bson:"repo_url,omitempty"`
	SystemURL          string `bson:"system_url,omitempty"`
	ProblemDescription string `bson:"problem_description,omitempty"`
}

// Session armazena o estado da conversa e outros dados do usuário.
type Session struct {
	UserID      string          `bson:"user_id"`
	State       string          `bson:"state"`
	Domain      string          `bson:"domain,omitempty"`
	PreAnalysis PreAnalysisData `bson:"pre_analysis,omitempty"`
}
