package routes_handle_functions

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"config"
	"db_adapter"
	"logging"
	"mycrypto"
	"structures"
)


func AuthHandlerOverride(
	handler func(
		response http.ResponseWriter,
		request *http.Request, 
		user_data structures.UserDataStructure)) http.HandlerFunc {

		return func(response http.ResponseWriter, request *http.Request) {
			var user_data structures.UserDataStructure

			cookie, err := request.Cookie(config.Token.Name)
			if err != nil {
				user_data.Username = ""
				user_data.Authentificated = false
				logging.Log.Printf("all_cookies: %v, Cookie: %s, Error: %v\n", request.Cookies(), cookie, err)
			} else {
				// Here should be token validation and getting username
				user_data.Username = "SomeUser"
				user_data.Authentificated = true
				logging.Log.Printf("Cookie: %s\n", cookie) // TODO Delete after cookie validation
			}

			handler(response, request, user_data)
	}
}

func HomePageHandler(
	response http.ResponseWriter, 
	request *http.Request, 
	user_data structures.UserDataStructure) {

	config.TemplatesParams.HomePage.UserData.Username = user_data.Username
	config.TemplatesParams.HomePage.UserData.Authentificated = user_data.Authentificated

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

func EmployeeLoginPageHandler(
	response http.ResponseWriter, 
	request *http.Request) {

    err := config.Templates.ExecuteTemplate(
		response, 
		config.RoutesHandlersInfo.EmployeeLogin.TemplateName, 
		config.TemplatesParams.EmployeeLogin) 

	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersInfo.EmployeeLogin.TemplateName, err)
		http.NotFound(response, request)
	}
}

func AuthEmployeeHandler(
	response http.ResponseWriter, 
	request *http.Request) {
		
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

	var username, password string = request.FormValue("username"), request.FormValue("password")

	db_adapter := db_adapter.DatabaseAdapter{}
	err := db_adapter.OpenConnection()
	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersInfo.SaveEmployee.TemplateName, err)
		http.Redirect(response, request, config.RoutesHandlersInfo.HomePage.TemplateName, http.StatusInternalServerError)
	}

	var encrypted_password string
	encrypted_password, err = mycrypto.Encrypt(password)
	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersInfo.SaveEmployee.TemplateName, err)
		http.Redirect(response, request, config.RoutesHandlersInfo.HomePage.TemplateName, http.StatusInternalServerError)
	}

	err = db_adapter.SaveEmployee(
		username, 
		encrypted_password)

	switch (err) {
		case nil:
			break
		case gorm.ErrRegistered:
			// TODO вынести в отдельный хэндлер, универсальный для всех типов ошибок + подумать над шаблонами для данных ошибок, чтобы юзер мог вернуться на главную страницу
			err_text := fmt.Sprintf("User with username=%s already registered. Error: %v\n", username, err) 
			logging.Log.Println(err_text)
			response.WriteHeader(http.StatusConflict)
			fmt.Fprintln(response, err_text)
			return
		default:
			logging.LogTemplateExecuteError(config.RoutesHandlersInfo.SaveEmployee.TemplateName, err)
			http.Redirect(
				response, 
				request, 
				config.RoutesHandlersInfo.HomePage.TemplateName, 
				http.StatusInternalServerError)
	}

	err = db_adapter.SaveToken(username, "sometokenname")
	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersInfo.SaveEmployee.TemplateName, err)
		http.Redirect(response, request, config.RoutesHandlersInfo.HomePage.TemplateName, http.StatusInternalServerError)
	}

	err = db_adapter.CloseConnection()
	if err != nil {
		logging.LogTemplateExecuteError(config.RoutesHandlersInfo.SaveEmployee.TemplateName, err)
		http.Redirect(response, request, config.RoutesHandlersInfo.HomePage.TemplateName, http.StatusInternalServerError)
	}

	// TODO need to add token generation here AND MOVE TO LOGIN PAGE. ALSO CHANGE REDIRECT TO LOGIN PAGE
	cookie1 := http.Cookie{
		Name: config.Token.Name, 
		Value: "SomeTokenValue", 
		MaxAge: config.Token.LifeTime,
		Path: config.Token.Path}
    http.SetCookie(response, &cookie1)

	// Используем прямой вызов хэндлера вместо http.Redirect из-за бесконечной переадресации в связи с тем, как строится URL при редиректе
	http.RedirectHandler(config.RoutesHandlersInfo.HomePage.URLPath, http.StatusFound).ServeHTTP(response, request)
}
