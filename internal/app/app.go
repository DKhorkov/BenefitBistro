package app

import (
	"config"
	"net/http"
	"routes_handle_functions"
)

func Run() {
	server := http.NewServeMux()

    // Подключаем обработку стилей для всех юрлов
    server.Handle(
        config.RoutesHandlersInfo.StaticFiles.URLPath, 
        http.StripPrefix(
            config.RoutesHandlersInfo.StaticFiles.URLPath, config.StaticFiles))

    // Создаем обработку юрлов
    server.HandleFunc(config.RoutesHandlersInfo.HomePage.URLPath, routes_handle_functions.HomePageHandler)
    server.HandleFunc(
        config.RoutesHandlersInfo.EmployeeRegister.URLPath, 
        routes_handle_functions.EmployeeRegisterPageHandler)

    http.ListenAndServe(config.ServerParameters.Host + ":" + config.ServerParameters.Port, server)
}
