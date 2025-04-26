package main

import (
	"fmt"
	"net/http"

	"flutter_package_syncer_server/config"
)

func main() {
	fmt.Println("Program is started!")
	println("Server started on http://localhost:8099")
	http.ListenAndServe(":8099", config.Routes())
}
