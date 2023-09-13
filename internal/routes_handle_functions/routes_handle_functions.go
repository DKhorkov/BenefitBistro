package routes_handle_functions

import (
	"config"
	"net/http"
	"structures"
)

func home_page_handler(response http.ResponseWriter, request *http.Request) {
	// Используем ExecuteTemplate метод, а не простой Execute, чтобы работало наследование шаблонов
    config.Templates.ExecuteTemplate(response, config.RoutesHandlersNames.HomePage, nil) 
}

var RouteHandlers structures.RouteHandlersStructure = structures.RouteHandlersStructure{
	HomePage: home_page_handler,
}
