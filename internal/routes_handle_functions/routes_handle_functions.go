package routes_handle_functions

import (
	"fmt"
	"html/template"
	"net/http"

	"gorm.io/gorm"

	"auth"
	"config"
	"db_adapter"
	"logging"
	"mycrypto"
	"structures"
)

// Чтобы использовать наследование HTML шаблонов в го, необходимо передать их все в ParseFiles
var templates = template.Must(template.ParseFiles(config.TemplateFiles...))


func AuthHandlerDelegat(
	handler func(
		response http.ResponseWriter,
		request *http.Request, 
		user_data structures.UserDataStructure),
	need_to_redirect bool) http.HandlerFunc {

		return func(response http.ResponseWriter, request *http.Request) {
			var user_data structures.UserDataStructure

			token, err := auth.GetAccessToken(request)
			if err != nil && need_to_redirect {
				http.Redirect(response, request, config.RoutesHandlersInfo.EmployeeLogin.URLPath, http.StatusSeeOther)
				return
			} else if err != nil {
				user_data.Username = ""
				user_data.Authentificated = false
			} else {
				db_adapter := db_adapter.DatabaseAdapter{}
				err := auth.ValidateToken(
					token, 
					need_to_redirect, 
					&user_data, 
					response, 
					request, 
					db_adapter)

				if err != nil {
					logging.Log.Printf("Error during AuthDelegat: %v\n", err)
					return
				}
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
		err := templates.ExecuteTemplate(
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

		err := templates.ExecuteTemplate(
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

		var username, password string = request.FormValue("username"), request.FormValue("password")
		hashed_password, err := mycrypto.HashPassword(password)
		if err != nil {
			logging.LogTemplateExecuteError(config.RoutesHandlersInfo.AuthEmployee.TemplateName, err)
			http.Redirect(
				response, 
				request, 
				config.RoutesHandlersInfo.EmployeeLogin.TemplateName, 
				http.StatusInternalServerError)

			return
		}
	
		db_adapter := db_adapter.DatabaseAdapter{}

		var corresponds bool
		corresponds, err = db_adapter.CompareEmployeeAuthData(username, hashed_password)
		if err != nil {
			// TODO создать отдельный метод для логирования и переадресации, чтобы не дублироваться
			logging.LogTemplateExecuteError(config.RoutesHandlersInfo.AuthEmployee.TemplateName, err)
			http.Redirect(
				response, 
				request, 
				config.RoutesHandlersInfo.EmployeeLogin.TemplateName, 
				http.StatusResetContent)

			return
		} else if !corresponds {
			logging.LogTemplateExecuteError(config.RoutesHandlersInfo.AuthEmployee.TemplateName, err)
			http.Redirect(
				response, 
				request, 
				config.RoutesHandlersInfo.EmployeeLogin.TemplateName, 
				http.StatusResetContent)

			return
		}

		// TODO need to change Token Generation to JWT
		token, err := mycrypto.Encrypt(fmt.Sprintf("%v+%v", username, password))
		if err != nil {
			logging.LogTemplateExecuteError(config.RoutesHandlersInfo.AuthEmployee.TemplateName, err)
			http.Redirect(
				response, 
				request, 
				config.RoutesHandlersInfo.EmployeeLogin.TemplateName, 
				http.StatusInternalServerError)
			
			return
		}

		token = config.Token.EmployeePrefix + token
		err = db_adapter.SaveEmployeeToken(username, token)
		if err != nil {
			logging.LogTemplateExecuteError(config.RoutesHandlersInfo.AuthEmployee.TemplateName, err)
			http.Redirect(
				response, 
				request, 
				config.RoutesHandlersInfo.EmployeeLogin.TemplateName, 
				http.StatusInternalServerError)
			
			return
		}

		cookie := &http.Cookie{
			Name: config.Token.Name, 
			Value: token, 
			MaxAge: config.Token.LifeTime,
			Path: config.Token.Path}
		http.SetCookie(response, cookie)

		// Используем прямой вызов хэндлера вместо http.Redirect из-за бесконечной переадресации в связи с тем, как строится URL при редиректе
		http.RedirectHandler(config.RoutesHandlersInfo.HomePage.URLPath, http.StatusFound).ServeHTTP(response, request)
}

func EmployeeRegisterPageHandler(
	response http.ResponseWriter, 
	request *http.Request) {

		err := templates.ExecuteTemplate(
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
		hashed_password, err := mycrypto.HashPassword(password)
		if err != nil {
			logging.LogTemplateExecuteError(config.RoutesHandlersInfo.SaveEmployee.TemplateName, err)
			http.Redirect(
				response, 
				request, 
				config.RoutesHandlersInfo.HomePage.TemplateName, 
				http.StatusInternalServerError)

			return
		}

		db_adapter := db_adapter.DatabaseAdapter{}
		err = db_adapter.SaveEmployee(
			username, 
			hashed_password)

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

				return
		}

		// Используем прямой вызов хэндлера вместо http.Redirect из-за бесконечной переадресации в связи с тем, как строится URL при редиректе
		http.RedirectHandler(
			config.RoutesHandlersInfo.EmployeeLogin.URLPath, 
			http.StatusSeeOther).ServeHTTP(response, request)
}

func LogoutHandler(
	response http.ResponseWriter, 
	request *http.Request) {

		token, err := auth.GetAccessToken(request)
		if err != nil {
			http.RedirectHandler(config.RoutesHandlersInfo.HomePage.URLPath, http.StatusFound).ServeHTTP(response, request)
			return
		}

		db_adapter := db_adapter.DatabaseAdapter{}
		db_adapter.DeleteToken(token)

		cookie := &http.Cookie{
			Name: config.Token.Name, 
			Value: "", 
			MaxAge: -1,
			Path: config.Token.Path}

		http.SetCookie(response, cookie)
		http.RedirectHandler(config.RoutesHandlersInfo.HomePage.URLPath, http.StatusFound).ServeHTTP(response, request)
}
