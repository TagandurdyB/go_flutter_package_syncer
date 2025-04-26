package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type API struct{}

func (api API) FlutterDoctor(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	localOutput, err := exec.Command("flutter", "doctor").CombinedOutput()
	if err != nil {
		http.Error(w, "Failed to run flutter doctor locally: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(string(localOutput))
}

func (api API) GetPaths(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, "Failed to get home directory", http.StatusInternalServerError)
		return
	}

	// Paths to scan
	scanDirs := []string{
		filepath.Join(homeDir, ".gradle"),
		filepath.Join(homeDir, ".pub-cache"),
	}

	var allPaths []string

	for _, dir := range scanDirs {
		// Make sure directory exists before walking
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue // Skip non-existent dirs
		}

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// Skip problematic files but log or handle if needed
				return nil
			}
			if !info.IsDir() {
				// Normalize and trim homeDir prefix for consistent paths
				relPath := strings.TrimPrefix(path, homeDir+string(os.PathSeparator))
				allPaths = append(allPaths, relPath)
			}
			return nil
		})

		if err != nil {
			http.Error(w, "Error walking path: "+dir+" - "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Save paths to a file
	outputFile := "server_paths.txt"
	f, err := os.Create(outputFile)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	for _, p := range allPaths {
		f.WriteString(p + "\n")
	}

	// Serve the file back to client
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\"server_paths.txt\"")

	http.ServeFile(w, r, outputFile)

}
