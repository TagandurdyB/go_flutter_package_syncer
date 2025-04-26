package main

import (
	"fmt"
	"net/http"

	"flutter_package_syncer/config"
)

func main() {
	fmt.Println("Program is started!")
	fmt.Println("You can acces from 127.0.0.1:8080")
	http.ListenAndServe(":8080", config.Routes())
}
