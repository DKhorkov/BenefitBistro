package app

import (
	"config"
	"net/http"
	"routes_handle_functions"
)

func Run() {
	server := http.NewServeMux()

    // Подключаем обработку стилей для всех юрлов
    server.Handle(config.URLPaths.StaticFiles, http.StripPrefix(config.URLPaths.StaticFiles, config.StaticFiles))

    // Создаем обработку юрлов
    server.HandleFunc(config.URLPaths.HomePage, routes_handle_functions.RouteHandlers.HomePage)

    http.ListenAndServe(config.ServerParameters.Host + ":" + config.ServerParameters.Port, server)
}
