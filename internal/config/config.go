package config

import (
	"time"

	"paths_and_folders"
	"structures"
)


var TemplateFiles []string = []string{
	"web/templates/navbar_core.html",
	"web/templates/non_authentificated_header.html",
	"web/templates/authentificated_header.html",
	"web/templates/footer.html",
	"web/templates/homepage.html",
	"web/templates/headTags.html",
	"web/templates/employeeRegister.html",
	"web/templates/employeeLogin.html",
}

var StaticFilesDir string = "web/static/"

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
	SaveEmployee: structures.RouteInfoStructure{
		TemplateName: "saveEmployee",
		URLPath: "/saveEmployee/",
	},
	EmployeeLogin: structures.RouteInfoStructure{
		TemplateName: "employeeLogin",
		URLPath: "/employeeLogin/",
	},
	AuthEmployee: structures.RouteInfoStructure{
		TemplateName: "authEmployee",
		URLPath: "/authEmployee/",
	},
	Logout: structures.RouteInfoStructure{
		TemplateName: "logout",
		URLPath: "/logout/",
	},
}

var ServerParameters structures.ServerParametersStructure = structures.ServerParametersStructure{
	Port: "8080",
}

var workspace_dir, _ = paths_and_folders.GetWorkspaceDir()
var temporatyFolder string = workspace_dir + "/tmp/"
var LogDir string = temporatyFolder + "logs/"
var LogPath string = LogDir + "log_file.log"

// Указатель, чтобы можно было динамически менять в процессе рабоыт приложения
var TemplatesParams *structures.TemplateParamsStructure = &structures.TemplateParamsStructure{
	HomePage: structures.TemplateDataStructure{
		PageName: "Главная страница",
	},
	EmployeeRegister: structures.TemplateDataStructure{
		PageName: "Регистрация сотрудника",
	},
}

var DatabaseFolder string = temporatyFolder + "database/"
var DatabaseName string = DatabaseFolder + "BenefitBistro.db"

// Следует хранить в окружении или скрыть в ином месте. Размер должен составлять 16, 24 или 32 байта.
var CryptKey []byte = []byte("574d93e6298df2e83e5c6b4dc63ae928") 

var Token structures.TokenStruct = structures.TokenStruct{
	Name: "Access-Token",
	LifeTime: 2 * 7 * 24 * 60 * 60, // 2 weeks
	Path: "/", // Path должен быть общим, чтобы кукис распространялись на весь сайт
	ExpiresDuration: time.Second * 2 * 7 * 24 * 60 * 60, // Expires in 2 weeks
	HirerPrefix: "hirer_",
	EmployeePrefix: "employee_",
}
