package controllers

import (
	"net/http"
	// "path/filepath"

	"github.com/julienschmidt/httprouter"
)

type FlutterControl struct{}

func (flutter FlutterControl) PubGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// repoName := params.ByName("repo_name")
	// repoPath := filepath.Join("storage/repos", repoName)
}

// func (flutter FlutterControl) BuildHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

// 	var req models.CloneRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	repoPath := filepath.Join("storage/repos", req.RepoURL)
// 	var cmd *exec.Cmd

// 	switch req.Action {
// 	case "pubget":
// 		cmd = exec.Command("flutter", "pub", "get")
// 	case "apk":
// 		cmd = exec.Command("flutter", "build", "apk")
// 	case "appbundle":
// 		cmd = exec.Command("flutter", "build", "appbundle")
// 	default:
// 		http.Error(w, "Unknown build action", http.StatusBadRequest)
// 		return
// 	}

// 	// Ensure we're in the repo directory and on the correct branch
// 	checkoutCmd := exec.Command("git", "-C", repoPath, "checkout", req.Branch)
// 	if out, err := checkoutCmd.CombinedOutput(); err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to checkout branch:\n%s", string(out)), http.StatusInternalServerError)
// 		return
// 	}

// 	// Set working directory for the Flutter command
// 	cmd.Dir = repoPath
// 	output, err := cmd.CombinedOutput()

// 	w.Header().Set("Content-Type", "application/json")
// 	resp := map[string]string{
// 		"output": string(output),
// 	}
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		resp["error"] = err.Error()
// 	}

// 	json.NewEncoder(w).Encode(resp)
// }
