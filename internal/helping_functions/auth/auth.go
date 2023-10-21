package auth

import (
	"fmt"
	"logging"
	"net/http"
	"strings"

	"config"
	"db_adapter"
	"structures"
)


func ValidateToken(
	token string, 
	need_to_redirect bool,
	user_data *structures.UserDataStructure, 
	response http.ResponseWriter, 
	request *http.Request,
	db_adapter db_adapter.DatabaseAdapter) error {

		if strings.HasPrefix(token, config.Token.EmployeePrefix) {
			user, err := db_adapter.ValidateEmployeeToken(token)
			return processTokenValidation(
				err, 
				need_to_redirect, 
				user_data, 
				user.Username, 
				config.RoutesHandlersInfo.EmployeeLogin.URLPath, 
				response, 
				request)
		} else if strings.HasPrefix(token, config.Token.HirerPrefix) {
			user, err := db_adapter.ValidateHirerToken(token)
			return processTokenValidation(
				err, 
				need_to_redirect, 
				user_data, 
				user.Username, 
				config.RoutesHandlersInfo.HirerLogin.URLPath, 
				response, 
				request)
		} else {
			return fmt.Errorf("invalid prefix for token: %v", token)
		}
}

func processTokenValidation(
	err error,
	need_to_redirect bool,
	user_data *structures.UserDataStructure,
	username string,
	redirectURL string,
	response http.ResponseWriter, 
	request *http.Request) error {

		if err != nil && need_to_redirect {
			logging.Log.Printf("Error during AuthDelegat: %v\n", err)
			http.Redirect(response, request, redirectURL, http.StatusSeeOther)
		} else if err != nil {
			logging.Log.Printf("Error during AuthDelegat: %v\n", err)
			user_data.Username = ""
			user_data.Authentificated = false
		} else {
			user_data.Username = username
			user_data.Authentificated = true
		}
		
		return err
}

func GetAccessToken(request *http.Request) (string, error) {
	cookie, err := request.Cookie(config.Token.Name)
	if err != nil {
		return "", err
	}

	token := strings.TrimPrefix(cookie.String(), fmt.Sprintf("%v=", config.Token.Name))
	return token, nil
}
