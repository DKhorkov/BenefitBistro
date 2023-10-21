package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"auth/testdata"
	"db_adapter"
	"paths_and_folders"
	"structures"
)

func TestGetAccessToken(test *testing.T) {
	request := httptest.NewRequest(http.MethodGet, testdata.TokenParams.Path, nil)
	cookie := &http.Cookie{
		Name: testdata.TokenParams.Name, 
		Value: testdata.EmployeeToken, 
		Path: testdata.TokenParams.Path}

	invalid_cookie := &http.Cookie{
		Name: testdata.InvalidTokenParams.Name, 
		Value: testdata.EmployeeToken, 
		Path: testdata.InvalidTokenParams.Path}
	
	request.AddCookie(invalid_cookie)
	token, err := GetAccessToken(request)
	assert.False(test, err == nil)
	assert.True(test, token == "")

	request.AddCookie(cookie)
	token, err = GetAccessToken(request)
	assert.True(test, err == nil)
	assert.True(test, token == testdata.EmployeeToken)
}

func TestProcessTokenValidation(test *testing.T) {
	user_data := &structures.UserDataStructure{}
	request := httptest.NewRequest(http.MethodGet, testdata.RedirectURL, nil)
	response := httptest.NewRecorder()

	err := processTokenValidation(
		errors.New("TestError"),
		testdata.NeedToRedirect,
		user_data,
		testdata.Username,
		testdata.RedirectURL,
		response,
		request,
	)
	assert.False(test, err == nil)
	assert.True(test, user_data.Username == "")
	assert.True(test, user_data.Authentificated == false)

	err = processTokenValidation(
		errors.New("TestError"),
		testdata.DoNotNeedToRedirect,
		user_data,
		testdata.Username,
		testdata.RedirectURL,
		response,
		request,
	)
	assert.False(test, err == nil)
	assert.True(test, user_data.Username == "")
	assert.True(test, user_data.Authentificated == false)

	err = processTokenValidation(
		nil,
		testdata.NeedToRedirect,
		user_data,
		testdata.Username,
		testdata.RedirectURL,
		response,
		request,
	)
	assert.True(test, err == nil)
	assert.True(test, user_data.Username == testdata.Username)
	assert.True(test, user_data.Authentificated == true)

	err = processTokenValidation(
		nil,
		testdata.DoNotNeedToRedirect,
		user_data,
		testdata.Username,
		testdata.RedirectURL,
		response,
		request,
	)
	assert.True(test, err == nil)
	assert.True(test, user_data.Username == testdata.Username)
	assert.True(test, user_data.Authentificated == true)
}

func TestValidateToken(test *testing.T) {
	// Preparing database for Validation Token
	db_adapter := db_adapter.DatabaseAdapter{
		DatabaseFolder: testdata.DatabaseFolder,
		DatabaseName: testdata.DatabaseName,
	}
	
	err := db_adapter.SaveEmployee(testdata.Username, testdata.Password)
	assert.True(test, err == nil)

	err = db_adapter.SaveEmployeeToken(testdata.Username, testdata.EmployeeToken)
	assert.True(test, err == nil)

	user_data := &structures.UserDataStructure{}
	request := httptest.NewRequest(http.MethodGet, testdata.RedirectURL, nil)
	response := httptest.NewRecorder()

	err = ValidateToken(
		testdata.EmployeeToken, 
		testdata.DoNotNeedToRedirect, 
		user_data, 
		response, 
		request, 
		db_adapter)
	assert.True(test, err == nil)

	err = ValidateToken(
		testdata.HirerToken, 
		testdata.DoNotNeedToRedirect, 
		user_data, 
		response, 
		request, 
		db_adapter)
	assert.True(test, err == nil)

	err = ValidateToken(
		testdata.RandomToken, 
		testdata.DoNotNeedToRedirect, 
		user_data, 
		response, 
		request, 
		db_adapter)
	assert.False(test, err == nil)

	err = paths_and_folders.DeletePath(testdata.DatabaseName)
	assert.True(test, err == nil)
}