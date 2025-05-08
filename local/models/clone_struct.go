package models

type CloneRequest struct {
	RepoURL string `json:"repoUrl"`
}

type GitResponse struct {
	Message  string   `json:"message"`
	Status   bool     `json:status`
	Error    string   `json:"error,omitempty"`
	Branches []string `json:"branches"`
}
