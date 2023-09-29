package routes_handle_functions

import (
	"net/http"

	"config"
	"db_adapter"
	"logging"
	"mycrypto"
)


func AuthHandlerOverride(
	handler func(
		response http.ResponseWriter,
		request *http.Request, 
		username string, 
		authentificated bool)) http.HandlerFunc {

		return func(response http.ResponseWriter, request *http.Request) {
			var authentificated bool
			var username string

			cookie, err := request.Cookie(config.Token.Name)
			if err != nil {
				authentificated = false
				username = ""
				logging.Log.Printf("all_cookies: %v, Cookie: %s, Error: %v\n", request.Cookies(), cookie, err)
			} else {
				// Here should be token validation and getting username
				authentificated = true
				username = "SomeUser"
				logging.Log.Printf("Cookie: %s\n", cookie) // TODO Delete after cookie validation
			}

			handler(response, request, username, authentificated)
	}
}

func HomePageHandler(
	response http.ResponseWriter, 
	request *http.Request, 
	username string, 
	authentificated bool) {

	config.TemplatesParams.HomePage.UserData.Username = username
	config.TemplatesParams.HomePage.UserData.Authentificated = authentificated

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

func EmployeeRegisterPageHandler(
	response http.ResponseWriter, 
	request *http.Request) {

    err := config.Templates.ExecuteTemplate(
		response, 
		config.RoutesHandlersInfo.EmployeeRegister.TemplateName, 
		config.TemplatesParams.EmployeeRegister) 

	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersInfo.EmployeeRegister.TemplateName, err)
		http.NotFound(response, request)
	}
}

func SaveEmployeeHandler(
	response http.ResponseWriter, 
	request *http.Request) {

	db_adapter := db_adapter.DatabaseAdapter{}
	err := db_adapter.OpenConnection()
	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersInfo.SaveEmployee.TemplateName, err)
		http.Redirect(response, request, config.RoutesHandlersInfo.HomePage.TemplateName, http.StatusInternalServerError)
	}

	var encrypted_password string
	encrypted_password, err = mycrypto.EncryptMessage(request.FormValue("password"))
	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersInfo.SaveEmployee.TemplateName, err)
		http.Redirect(response, request, config.RoutesHandlersInfo.HomePage.TemplateName, http.StatusInternalServerError)
	}

	err = db_adapter.SaveEmployee(
		request.FormValue("username"), 
		encrypted_password)

	// TODO Сделать проверку на тип ошибки и в зависимости от ошибки возвращать разные статусы и сообщения на фронт
	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersInfo.SaveEmployee.TemplateName, err)
		http.Redirect(response, request, config.RoutesHandlersInfo.HomePage.TemplateName, http.StatusInternalServerError)
	}

	db_adapter.CloseConnection()

	// TODO need to add token generation here AND MOVE TO LOGIN PAGE. ALSO CHANGE REDIRECT TO LOGIN PAGE
	cookie1 := http.Cookie{
		Name: config.Token.Name, 
		Value: "SomeTokenValue", 
		MaxAge: config.Token.LifeTime,
		Path: "/"} // Path должен быть общим, чтобы кукис распространялись на весь сайт
    http.SetCookie(response, &cookie1)

	// Используем прямой вызов хэндлера вместо http.Redirect из-за бесконечной переадресации в связи с тем, как строится URL при редиректе
	http.RedirectHandler(config.RoutesHandlersInfo.HomePage.URLPath, http.StatusFound).ServeHTTP(response, request)
}
