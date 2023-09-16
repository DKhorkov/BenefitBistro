package routes_handle_functions

import (
	"config"
	"logging"
	"net/http"
	"structures"
)

func home_page_handler(response http.ResponseWriter, request *http.Request) {
	// Используем ExecuteTemplate метод, а не простой Execute, чтобы работало наследование шаблонов
    err := config.Templates.ExecuteTemplate(response, "bla", nil) 

	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersNames.HomePage, err)
		http.NotFound(response, request)
	}
}

var RouteHandlers structures.RouteHandlersStructure = structures.RouteHandlersStructure{
	HomePage: home_page_handler,
}
