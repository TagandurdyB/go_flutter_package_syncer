package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"

	"flutter_package_syncer/helpers"
	"flutter_package_syncer/models"
)

type GitControl struct{}

func (git GitControl) Clone(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Parse JSON body
	var req models.CloneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	repoURL := req.RepoURL
	if repoURL == "" {
		http.Error(w, "Repository URL is required", http.StatusBadRequest)
		return
	}

	repoName, output, status := helpers.GitService{}.Clone(repoURL)

	if !status {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.GitResponse{
			Message: "Failed to clone repository",
			Error:   string(output),
		})
	}

	json.NewEncoder(w).Encode(models.GitResponse{
		Message: fmt.Sprintf("Repository '%s' cloned successfully!", repoName),
		Status:  true,
	})

}

func (git GitControl) Branches(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	repoName := params.ByName("repo_name")
	repoPath := filepath.Join("storage/repos", repoName)

	output, branches, status := helpers.GitService{}.Branches(repoPath)

	w.Header().Set("Content-Type", "application/json")
	if !status {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.GitResponse{
			Message: "Failed to clone repository",
			Error:   string(output),
			Status:  false,
		})
	} else {
		json.NewEncoder(w).Encode(models.GitResponse{
			Message:  "Repository '%s' branches get successfully!",
			Branches: branches,
			Status:   true,
		})
	}
}

func (git GitControl) Pull(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	repo := r.URL.Query().Get("repo")
	branch := r.URL.Query().Get("branch")

	if repo == "" || branch == "" {
		http.Error(w, "Missing repo or branch parameter", http.StatusBadRequest)
		return
	}

	// Assuming your repos are under ./repos/
	repoPath := fmt.Sprintf("./storage/repos/%s", repo)

	// Prepare Git command
	output, status := helpers.GitService{}.Pull(repoPath, branch)

	w.Header().Set("Content-Type", "application/json")
	if !status {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.GitResponse{
			Message: string(fmt.Sprintf("Git pull failed!")),
			Error:   string(output),
			Status:  false,
		})
		// http.Error(w, fmt.Sprintf("Git pull failed:\n%s\nError: %s", output, err), http.StatusInternalServerError)

	} else {
		w.Write(output)
	}
}
