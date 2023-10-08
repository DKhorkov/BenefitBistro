package app

import (
	"net/http"

	"config"
	"routes_handle_functions"
)

func Run() {
	server := http.NewServeMux()

    // Подключаем обработку стилей для всех юрлов
    var staticFiles = http.FileServer(http.Dir(config.StaticFilesDir))
    server.Handle(
        config.RoutesHandlersInfo.StaticFiles.URLPath, 
        http.StripPrefix(
            config.RoutesHandlersInfo.StaticFiles.URLPath, staticFiles))

    // Создаем обработку юрлов
    server.HandleFunc(
        config.RoutesHandlersInfo.HomePage.URLPath,
        routes_handle_functions.AuthHandlerDelegat(
            routes_handle_functions.HomePageHandler,
            false))
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
    server.HandleFunc(
        config.RoutesHandlersInfo.Logout.URLPath, 
        routes_handle_functions.LogoutHandler)

    http.ListenAndServe(config.ServerParameters.Host + ":" + config.ServerParameters.Port, server)
}
