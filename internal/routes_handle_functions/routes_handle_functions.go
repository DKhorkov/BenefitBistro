package routes_handle_functions

import (
	"config"
	"db_adapter"
	"logging"
	"net/http"
)

func HomePageHandler(response http.ResponseWriter, request *http.Request) {
	// Используем ExecuteTemplate метод, а не простой Execute, чтобы работало наследование шаблонов
    err := config.Templates.ExecuteTemplate(
		response, 
		config.RoutesHandlersInfo.HomePage.TemplateName, 
		config.TemplatesParams.HomePage) 

	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersInfo.HomePage.TemplateName, err)
		http.NotFound(response, request)
	}
}

func EmployeeRegisterPageHandler(response http.ResponseWriter, request *http.Request) {
	// Используем ExecuteTemplate метод, а не простой Execute, чтобы работало наследование шаблонов
    err := config.Templates.ExecuteTemplate(
		response, 
		config.RoutesHandlersInfo.EmployeeRegister.TemplateName, 
		config.TemplatesParams.EmployeeRegister) 

	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersInfo.EmployeeRegister.TemplateName, err)
		http.NotFound(response, request)
	}
}

func SaveEmployeeHandler(response http.ResponseWriter, request *http.Request) {
	logging.Log.Printf("saving employee %v\n", request)
	db_adapter := db_adapter.DatabaseAdapter{}
	err := db_adapter.OpenConnection()
	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersInfo.SaveEmployee.TemplateName, err)
		http.Redirect(response, request, config.RoutesHandlersInfo.HomePage.TemplateName, http.StatusInternalServerError)
	}

	err = db_adapter.SaveEmployee(
		request.FormValue("username"), 
		request.FormValue("password"))

	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersInfo.SaveEmployee.TemplateName, err)
		http.Redirect(response, request, config.RoutesHandlersInfo.HomePage.TemplateName, http.StatusInternalServerError)
	}

	db_adapter.CloseConnection()

	// Используем прямой вызов хэндлера вместо http.Redirect из-за бесконечной переадресации в связи с тем, как строится URL при редиректе
	http.RedirectHandler(config.RoutesHandlersInfo.HomePage.URLPath, http.StatusFound).ServeHTTP(response, request)
}
