package routes_handle_functions

import (
	"config"
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