package config

import (
	"github.com/julienschmidt/httprouter"

	controllers "flutter_package_syncer_server/controllers"
)

func Routes() *httprouter.Router {
	r := httprouter.New()
	// API Endpoints
	r.GET("/api/flutter-doctor", controllers.API{}.FlutterDoctor)
	r.GET("/api/get-paths", controllers.API{}.GetPaths)
	//Serve Files
	// r.ServeFiles("/views/assets/*filepath", http.Dir("views/assets/"))
	return r
}
