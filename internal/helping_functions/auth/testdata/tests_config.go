package testdata

import (
	"paths_and_folders"
	"structures"
)
	
var workspace_dir, _ = paths_and_folders.GetWorkspaceDir()

var (
	TokenParams structures.TokenStruct = structures.TokenStruct{
		Name: "Access-Token",
		Path: "/",
	}

	InvalidTokenParams structures.TokenStruct = structures.TokenStruct{
		Name: "Invalid-Token",
		Path: "/",
	}

	EmployeeToken string = "employee_token"
	HirerToken string = "hirer_token"
	RandomToken string = "random_token"
	Username string = "TestUsername"
	Password string = "TestPassword"
	RedirectURL string = "/"
	NeedToRedirect bool = true
	DoNotNeedToRedirect bool = false

	testsTemporatyFolder string = workspace_dir + "/tmp/"
	DatabaseFolder string = testsTemporatyFolder + "tests_folder/"
	DatabaseName string = DatabaseFolder + "test_database.db"
)