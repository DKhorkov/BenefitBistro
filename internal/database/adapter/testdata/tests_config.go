package testdata

import (
	"paths_and_folders"
)

var workspace_dir, _ = paths_and_folders.GetWorkspaceDir()

var (
	testsTemporatyFolder string = workspace_dir + "/tmp/"
	DatabaseFolder string = testsTemporatyFolder + "tests_folder/"
	DatabaseName string = DatabaseFolder + "test_database.db"

	TableName string = "employee_users"

	Username string = "test_user"
	Password string = "test_password"
	InvalidPassword string = "some_invalid_password"
	EmployeeToken string = "employee_test_token"
	HirerToken string = "hirer_test_token"
)
