package helpers

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	LocalHost    string
	ServerHost   string
	LocalPort    string
	ServerPort   string
	LocalDomain  string
	ServerDomain string
)

func InitEnv() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Load values from environment variables
	LocalHost = os.Getenv("LOCAL_HOST")
	ServerHost = os.Getenv("SERVER_HOST")
	LocalPort = os.Getenv("LOCAL_PORT")
	ServerPort = os.Getenv("SERVER_PORT")

	// Check if any variable is missing and print a warning
	if LocalHost == "" || ServerHost == "" || LocalPort == "" || ServerPort == "" {
		fmt.Println("Warning: One or more environment variables are not set correctly.")
	}

	LocalDomain = LocalHost + ":" + LocalPort
	ServerDomain = ServerHost + ":" + ServerPort

	// Optionally, you can print the loaded values to check
	// fmt.Println("Loaded Config:")
	// fmt.Println("LOCAL_HOST:", LocalHost)
	// fmt.Println("SERVER_HOST:", ServerHost)
	// fmt.Println("LOCAL_PORT:", LocalPort)
	// fmt.Println("SERVER_PORT:", ServerPort)
}
