package helpers

import (
	"path/filepath"
)

func Include(path string) []string {
	files, err := filepath.Glob("views/templates/*.html")
	ErrH("Error in Include: ", err)

	path_files, err := filepath.Glob("views/" + path + "/*.html")
	ErrH("Error in Include: ", err)

	files = append(files, path_files...)

	// fmt.Println("INCLUDS=========================================================")
	// for _, file := range files {
	// 	fmt.Println(file)
	// }
	// fmt.Println("/INCLUDS=========================================================")

	return files
}
