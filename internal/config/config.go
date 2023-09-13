package config

import (
	"html/template"
	"net/http"
	"structures"
)

var template_files []string = []string{
	"web/templates/base.html",
	"web/templates/header.html",
}


var static_files_dir string = "web/static/"
var StaticFiles = http.FileServer(http.Dir(static_files_dir))

// Чтобы использовать наследование HTML шаблонов в го, необходимо передать их все в ParseFiles
var Templates = template.Must(template.ParseFiles(template_files...))

var RoutesHandlersNames structures.RouteHandlersNamesStructure = structures.RouteHandlersNamesStructure{
	HomePage: "base",
}

var ServerParameters structures.ServerParametersStructure = structures.ServerParametersStructure{
	Port: "8010",
}

var URLPaths structures.URLPathsStructure = structures.URLPathsStructure{
	HomePage: "/",
	StaticFiles: "/static/",
}
