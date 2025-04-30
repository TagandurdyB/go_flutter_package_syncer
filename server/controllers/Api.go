package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"

	"flutter_package_syncer_server/helpers"
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

func (api API) UploadHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	r.ParseMultipartForm(32 << 20) // 32MB memory buffer

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempDir := "tmp"
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		http.Error(w, "Failed to create temp directory", http.StatusInternalServerError)
		return
	}
	dst, err := os.Create(tempDir + "/" + handler.Filename)
	if err != nil {
		http.Error(w, "Failed to create file on server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File %s uploaded successfully", handler.Filename)
}

func (api API) UnpackArchive(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	archivePath := "./tmp/packages.tar" // Make sure this matches your real file
	destDir := "./tmp/unpacked"

	fmt.Println("⏳ Starting unpack...")

	// Call the unpack helper function
	err := helpers.ExtractTarGz(archivePath, destDir)
	if err != nil {
		http.Error(w, "❌ Unpack failed: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("❌ Error during unpack:", err)
		return
	}

	fmt.Println("✅ Archive unpacked successfully")

	// Prepare response with a success message
	response := map[string]string{
		"message": "✅ Archive unpacked successfully",
	}

	// Set response header to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Convert response to JSON and send it in the body
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "❌ Error while marshalling JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the JSON response
	w.Write(jsonResponse)
}

func (api API) Sync(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Start of the sync process
	println("Sync process started")

	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, fmt.Sprintf("❌ Failed to get home directory: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Define the source and destination directories
	srcDir := "./tmp/unpacked/"
	destDir := homeDir + "/"

	// Print the source and destination paths for debugging
	fmt.Printf("Syncing from: %s to: %s\n", srcDir, destDir)

	// Sync files from ./tmp/unpacked to ~/*
	err = helpers.SyncFiles(srcDir, destDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("❌ Sync failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("✅ Sync completed successfully")
}

