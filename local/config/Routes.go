package config

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	controllers "flutter_package_syncer/controllers"
)

func Routes() *httprouter.Router {
	r := httprouter.New()
	// View
	r.GET("/", controllers.Dashboard{}.Index)
	// API Endpoints
	r.GET("/api/flutter-doctor", controllers.Dashboard{}.FlutterDoctor)
	r.GET("/api/package-diff", controllers.Dashboard{}.PackagesDiff)
	r.GET("/api/archive", controllers.Dashboard{}.Archive)
	r.POST("/api/upload", controllers.Dashboard{}.UploadPackages)
	r.GET("/api/sync-packages", controllers.Dashboard{}.SyncPackages)
	//
	//Serve Files
	r.ServeFiles("/views/assets/*filepath", http.Dir("views/assets/"))
	return r
}
