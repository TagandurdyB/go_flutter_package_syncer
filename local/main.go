package main

import (
	"fmt"
	"net/http"

	"flutter_package_syncer/config"
	"flutter_package_syncer/helpers"
)

func main() {
	helpers.InitEnv()
	fmt.Println("Program is started!")
	fmt.Println("You can acces from " + helpers.LocalDomain)
	http.ListenAndServe(helpers.LocalDomain, config.Routes())
}
