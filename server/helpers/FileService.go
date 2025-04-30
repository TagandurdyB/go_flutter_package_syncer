package helpers

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func MkDir(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		e := os.MkdirAll(path, 0755)
		return !ErrH("Err in mkdir(", path, ")", e)
	}
	return false
}

func IsExist(fileName string) bool {
	f, err := os.Open(fileName)
	f.Close()
	if err == nil {
		//logSave(fileName + " exist! (true)")
		return true
	} else if os.IsNotExist(err) {
		//logSave(fileName + " not exist! (false)")
		return false
	} else {
		logSave("Unknown error isExist(", fileName, "): ", err)
		return false
	}
}

func CreateFile(fileName string) bool {
	_, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logSave("Unknown error createFile(", fileName, "): ", err)
	} else if !IsExist(fileName) {
		_, err = os.Create(fileName)
		ErrH("Error createFile(", fileName, "): ", err)
		return true
	}
	// logSave("The file has already been created! : " + fileName)
	return false
}

func ReadFile(fileName string) (text []string) {
	file, _ := os.Open(fileName)
	if IsExist(fileName) {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			text = append(text, line)
		}
	} else {
		logSave("readFile : ", file, " not exist!")
	}
	file.Close()
	if len(text) == 0 {
		text = append(text, "")
	}
	return
}

func WriteFile(fileName string, text string) {
	file, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if IsExist(fileName) {
		_, err := file.WriteString(text)
		ErrH("Error WriteFile(", fileName, "): ", err)
	}
	defer file.Close()
}

func WriteJson(fileName string, jsonData []byte) {
	err := os.WriteFile(fileName, jsonData, fs.FileMode(0644))
	ErrH("Error WriteJson(", fileName, "):", err)
}

func AppendJson(fileName string, jsonData []byte) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	ErrH("Error AppendJson(", fileName, "):", err)
	defer file.Close()
	// var result []interface{}
	stat, err := file.Stat()
	ErrH("Error AppendJson(", fileName, "):", err)
	isEmpty := stat.Size() == 0
	if !isEmpty {
		jsonData = append([]byte{','}, jsonData...)
	}
	_, err = file.Write(jsonData)
	ErrH("Error AppendJson(", fileName, "):", err)

}

func ReadAllJson(fileName string) []byte {

	jsonData, err := os.ReadFile(fileName)
	ErrH("Error WriteJson(", fileName, "):", err)
	// err = json.Unmarshal(jsonData, &result)
	return jsonData
}

func AppendFile(fileName string, text string) {
	file, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if IsExist(fileName) {
		_, err := file.WriteString(text)
		ErrH("Error appendFile(", fileName, "): ", err)
	}
	defer file.Close()
}

func DeleteFile(fileName string) {
	err := os.Remove(fileName)
	ErrH("Error deleteFile(", fileName, "): ", err)
}

func SyncFiles(srcDir, destDir string) error {
	// Open the source directory
	sourceDir, err := os.Open(srcDir)
	if err != nil {
		return fmt.Errorf("failed to open source directory: %v", err)
	}
	defer sourceDir.Close()

	// Read the content of the source directory
	files, err := sourceDir.Readdir(-1) // -1 to read all files
	if err != nil {
		return fmt.Errorf("failed to read source directory: %v", err)
	}

	// Iterate over each file in the source directory
	for _, file := range files {
		// Construct full file paths
		srcFilePath := filepath.Join(srcDir, file.Name())
		destFilePath := filepath.Join(destDir, file.Name())

		// If it's a directory, recursively sync the contents
		if file.IsDir() {
			err := os.MkdirAll(destFilePath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %v", destFilePath, err)
			}
			// Recursively sync the files in the subdirectory
			err = SyncFiles(srcFilePath, destFilePath)
			if err != nil {
				return err
			}
		} else {
			// It's a file, copy it to the destination
			err := CopyFile(srcFilePath, destFilePath)
			if err != nil {
				return fmt.Errorf("failed to copy file %s to %s: %v", srcFilePath, destFilePath, err)
			}
		}
	}

	return nil
}

// copyFile copies a single file from src to dest
func CopyFile(src, dest string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	// Create the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	// Copy the contents from source to destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	// Return success
	return nil
}
