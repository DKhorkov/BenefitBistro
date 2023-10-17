package testdata

import (
	"paths_and_folders"
)

var workspace_dir, _ = paths_and_folders.GetWorkspaceDir()

var (
	testsTemporatyFolder string = workspace_dir + "/tmp/"
	DatabaseFolder string = testsTemporatyFolder + "tests_folder/"
	DatabaseName string = DatabaseFolder + "test_database.db"

	Username string = "test_user"
	Password string = "test_password"
	Token string = "test_token"
)
