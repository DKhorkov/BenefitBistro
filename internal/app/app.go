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
    server.HandleFunc(
        config.RoutesHandlersInfo.HomePage.URLPath,
        routes_handle_functions.AuthHandlerDelegat(
            routes_handle_functions.HomePageHandler,
            true))
    server.HandleFunc(
        config.RoutesHandlersInfo.EmployeeRegister.URLPath, 
        routes_handle_functions.EmployeeRegisterPageHandler)
    server.HandleFunc(
        config.RoutesHandlersInfo.SaveEmployee.URLPath, 
        routes_handle_functions.SaveEmployeeHandler)
    server.HandleFunc(
        config.RoutesHandlersInfo.EmployeeLogin.URLPath, 
        routes_handle_functions.EmployeeLoginPageHandler)
    server.HandleFunc(
        config.RoutesHandlersInfo.AuthEmployee.URLPath, 
        routes_handle_functions.AuthEmployeeHandler)

    http.ListenAndServe(config.ServerParameters.Host + ":" + config.ServerParameters.Port, server)
}
