package config

import (
	"html/template"
	"net/http"
	"structures"
)


var template_files []string = []string{
	"web/templates/header.html",
	"web/templates/footer.html",
	"web/templates/homepage.html",
	"web/templates/headTags.html",
	"web/templates/employeeRegister.html",
}

var static_files_dir string = "web/static/"
var StaticFiles = http.FileServer(http.Dir(static_files_dir))

// Чтобы использовать наследование HTML шаблонов в го, необходимо передать их все в ParseFiles
var Templates = template.Must(template.ParseFiles(template_files...))

var RoutesHandlersInfo structures.RouteHandlersInfoStructure = structures.RouteHandlersInfoStructure{
	StaticFiles: structures.RouteInfoStructure{
		URLPath: "/static/",
	},
	HomePage: structures.RouteInfoStructure{
		TemplateName: "homepage",
		URLPath: "/",
	},
	EmployeeRegister: structures.RouteInfoStructure{
		TemplateName: "employeeRegister",
		URLPath: "/employeeRegister/",
	},
}

var ServerParameters structures.ServerParametersStructure = structures.ServerParametersStructure{
	Port: "8090",
}

var LogDir string = "tmp/logs/"
var LogPath string = LogDir + "log_file"

var TemplatesParams structures.TemplateParamsStructure = structures.TemplateParamsStructure{
	HomePage: structures.TemplateData{
		PageName: "Главная страница",
	},
	EmployeeRegister: structures.TemplateData{
		PageName: "Регистрация сотрудника",
	},
}
