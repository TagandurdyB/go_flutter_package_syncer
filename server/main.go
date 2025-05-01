package main

import (
	"fmt"
	"net/http"

	"flutter_package_syncer_server/config"
	"flutter_package_syncer_server/helpers"
)

func main() {
	helpers.InitEnv()
	fmt.Println("Program is started!")
	println("Server started on http://"+helpers.Domain)
	http.ListenAndServe(helpers.Domain, config.Routes())
}
