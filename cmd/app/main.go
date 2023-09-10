package main

import (
	"html/template"
	"net/http"
)


var template_files []string = []string{
	"web/templates/base.html",
	"web/templates/header.html",
}

// Чтобы использовать наследование HTML шаблонов в го, необходимо передать их все в ParseFiles
var templates = template.Must(template.ParseFiles(template_files...))

func baseHandler(response http.ResponseWriter, request *http.Request) {
	// Используем ExecuteTemplate метод, а не простой Execute, чтобы работало наследование шаблонов
    templates.ExecuteTemplate(response, "base", nil) 
}


func main() {
    port := "8010"

    server := http.NewServeMux()

    // Подключаем обработку стилей для всех юрлов
    static_files := http.FileServer(http.Dir("web/static/"))
    server.Handle("/static/", http.StripPrefix("/static/", static_files))

    // Создаем обработку юрлов
    server.HandleFunc("/", baseHandler)

    http.ListenAndServe(":" + port, server)
}