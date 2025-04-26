package controllers

import (
	"bufio"
	"encoding/json"
	"html/template"
	"io"
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
	resp, err := http.Get("http://localhost:8099/api/flutter-doctor")
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
		"./diff",
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
	resp, err := http.Get("http://localhost:8099/api/get-paths")
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
