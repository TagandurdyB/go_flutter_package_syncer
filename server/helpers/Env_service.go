package helpers

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	Host   string
	Port   string
	Domain string
)

func InitEnv() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Load values from environment variables
	Host = os.Getenv("HOST")
	Port = os.Getenv("PORT")

	// Check if any variable is missing and print a warning
	if Host == "" || Port == ""  {
		fmt.Println("Warning: One or more environment variables are not set correctly.")
	}
	Domain = Host + ":" + Port

	// Optionally, you can print the loaded values to check
	// fmt.Println("Loaded Config:")
	// fmt.Println("Host:", Host)
	// fmt.Println("Port:", Port)
}
