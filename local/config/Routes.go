package config

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	controllers "flutter_package_syncer/controllers"
)

func Routes() *httprouter.Router {
	r := httprouter.New()
	dashboard := controllers.Dashboard{}
	git := controllers.GitControl{}
	flutter := controllers.FlutterControl{}
	// View
	r.GET("/", dashboard.Index)
	r.GET("/repos", dashboard.Repos)
	r.GET("/repos/:repo_name/:branch", dashboard.Branchs)
	// API Endpoints
	r.GET("/api/flutter-doctor", dashboard.FlutterDoctor)
	r.GET("/api/package-diff", dashboard.PackagesDiff)
	r.GET("/api/archive", dashboard.Archive)
	r.POST("/api/upload", dashboard.UploadPackages)
	r.GET("/api/sync-packages", dashboard.SyncPackages)
	//Git
	r.POST("/api/clone", git.Clone)
	r.GET("/api/branches/:repo_name", git.Branches)
	r.GET("/api/git-pull", git.Pull)
	//Flutter
	r.POST("/api/flutter/pub-get/:repo_name", flutter.PubGet)

	//
	//Serve Files
	r.ServeFiles("/views/assets/*filepath", http.Dir("views/assets/"))
	return r
}
