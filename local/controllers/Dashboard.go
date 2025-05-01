package controllers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"flutter_package_syncer/helpers"
)

type Dashboard struct{}

type ResultMessage struct {
	Local       string   `json:"local"`
	Server      string   `json:"server"`
	DiffMessage string   `json:"diff_message"`
	Diff        []string `json:"diff"`
}

type DiffResponse struct {
	Message string `json:"message"`
}

func (dashboard Dashboard) Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	temp := helpers.Include("dashboard")

	view, err := template.ParseFiles(temp...)

	helpers.ErrH("Error in Dashboard Index: ", err)
	data := make(map[string]interface{})
	// data["Students"] = models.Person{}.ReadAll()

	view.ExecuteTemplate(w, "Index", data)
}

func (dashboard Dashboard) FlutterDoctor(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	localOutput, err := exec.Command("flutter", "doctor").CombinedOutput()
	if err != nil {
		http.Error(w, "Failed to run flutter doctor locally: "+err.Error(), http.StatusInternalServerError)
		return
	}
	result := ResultMessage{
		Local: string(localOutput),
	}

	// TODO: Replace with real remote/server logic
	serverOutput, err := serverFlutterDoctor()
	if err != nil {
		result = ResultMessage{
			Local:  string(localOutput),
			Server: err.Error(),
		}
	} else {
		result = ResultMessage{
			Local:  string(localOutput),
			Server: serverOutput,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func serverFlutterDoctor() (string, error) {
	// Send request
	resp, err := http.Get("http://" + helpers.ServerDomain + "/api/flutter-doctor")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	return bodyString, nil
}

func (dashboard Dashboard) PackagesDiff(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, "Failed to get home directory", http.StatusInternalServerError)
		return
	}

	// Paths to scan
	scanDirs := []string{
		filepath.Join(homeDir, ".gradle"),
		filepath.Join(homeDir, ".pub-cache"),
		filepath.Join(homeDir, ".diff"),
	}

	var localPaths []string

	for _, dir := range scanDirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // skip errors silently
			}
			if !info.IsDir() {
				relPath := strings.TrimPrefix(path, homeDir+"/")
				localPaths = append(localPaths, relPath)
			}
			return nil
		})
		if err != nil {
			http.Error(w, "Error walking path: "+dir, http.StatusInternalServerError)
			return
		}
	}

	// Save to paths.txt
	err = os.WriteFile("paths.txt", []byte(strings.Join(localPaths, "\n")), 0644)
	if err != nil {
		http.Error(w, "Failed to write paths.txt", http.StatusInternalServerError)
		return
	}

	serverPaths, err := downloadServerPaths()
	if err != nil {
		json.NewEncoder(w).Encode(ResultMessage{
			Local:       "Listed 234 local paths successful!",
			Server:      "Server paths can't fetched!",
			DiffMessage: "Server paths can't fetched!",
			Diff:        []string{},
		})
	}

	onlyInLocal := []string{}
	onlyInLocal, _ = findDifferences(localPaths, serverPaths)
	err = os.WriteFile("only_in_local.txt", []byte(strings.Join(onlyInLocal, "\n")), 0644)
	if err != nil {
		json.NewEncoder(w).Encode(ResultMessage{
			Local:       "Listed " + strconv.Itoa(len(localPaths)) + " local paths successful!",
			Server:      "Listed " + strconv.Itoa(len(serverPaths)) + " server paths successful!",
			DiffMessage: "Failed to write differences!",
			Diff:        onlyInLocal,
		})
	} else {
		json.NewEncoder(w).Encode(ResultMessage{
			Local:       "Listed " + strconv.Itoa(len(localPaths)) + " local paths successful!",
			Server:      "Listed " + strconv.Itoa(len(serverPaths)) + " server paths successful!",
			DiffMessage: strconv.Itoa(len(onlyInLocal)) + " files found that do not exist on the server!",
			Diff:        onlyInLocal,
		})
	}

}

func downloadServerPaths() ([]string, error) {
	// Send request
	resp, err := http.Get("http://" + helpers.ServerDomain + "/api/get-paths")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "server_paths_*.txt")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name()) // clean up after reading
	defer tempFile.Close()

	// Write response into temp file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return nil, err
	}

	// Now open the file to read lines
	file, err := os.Open(tempFile.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func findDifferences(localPaths, serverPaths []string) (onlyInLocal, onlyInServer []string) {
	serverMap := make(map[string]bool)
	localMap := make(map[string]bool)

	for _, path := range serverPaths {
		serverMap[path] = true
	}

	for _, path := range localPaths {
		localMap[path] = true
		if !serverMap[path] {
			onlyInLocal = append(onlyInLocal, path)
		}
	}

	for _, path := range serverPaths {
		if !localMap[path] {
			onlyInServer = append(onlyInServer, path)
		}
	}

	return
}

func (dashboard Dashboard) Archive(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, "Failed to get home directory", http.StatusInternalServerError)
		return
	}

	// Read only_in_local.txt
	data, err := os.ReadFile("only_in_local.txt")
	if err != nil {
		http.Error(w, "Failed to read only_in_local.txt", http.StatusInternalServerError)
		return
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	// Create temp dir under current directory
	tempDir := "./archive-temp"
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		http.Error(w, "Failed to create temp directory", http.StatusInternalServerError)
		return
	}
	// defer os.RemoveAll(tempDir) // clean up after

	// Copy each file
	for _, relPath := range lines {
		relPath = strings.TrimSpace(relPath)
		if relPath == "" {
			continue
		}

		sourcePath := filepath.Join(homeDir, relPath)
		destPath := filepath.Join(tempDir, relPath)

		// Ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			http.Error(w, "Failed to create destination directory: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy the file
		err := helpers.CopyFile(sourcePath, destPath)
		if err != nil {
			http.Error(w, "Failed to copy file: "+sourcePath+" -> "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Now create the tar archive
	tarPath := "packages.tar"
	err = helpers.CompressToTarGz(tempDir, tarPath)
	if err != nil {
		http.Error(w, "Failed to create tar archive: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove temp directory
	_ = os.RemoveAll(tempDir)

	// Get file info for size and name
	info, err := os.Stat(tarPath)
	if err != nil {
		http.Error(w, "Failed to stat tar file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with JSON info
	response := map[string]interface{}{
		"name": info.Name(),
		"size": info.Size(),
		"path": tarPath,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (dashboard Dashboard) UploadPackages(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	filePath := "packages.tar"
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	// Prepare multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		http.Error(w, "Failed to create form file", http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(part, file); err != nil {
		http.Error(w, "Failed to copy file data", http.StatusInternalServerError)
		return
	}

	writer.Close()

	// Send POST request to server
	req, err := http.NewRequest("POST", "http://"+helpers.ServerDomain+"/api/upload", body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (dashboard Dashboard) SyncPackages(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// First request: Send to unpack API
	resp, err := http.Post("http://"+helpers.ServerDomain+"/api/unpack", "application/json", nil)
	if err != nil {
		http.Error(w, "❌ Failed to sync with server: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("❌ Server failed to unpack: %s", string(bodyBytes)), resp.StatusCode)
		return
	}

	// Second request: Send to sync API (or any other endpoint)
	resp2, err := http.Post("http://"+helpers.ServerDomain+"/api/sync", "application/json", nil)
	if err != nil {
		http.Error(w, "❌ Failed to sync with server: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp2.Body)
		http.Error(w, fmt.Sprintf("❌ Server failed to sync: %s", string(bodyBytes)), resp2.StatusCode)
		return
	}

	// ✅ Forward success response from the second API (sync API) to frontend
	w.WriteHeader(http.StatusOK)
	io.Copy(w, resp2.Body)
}
