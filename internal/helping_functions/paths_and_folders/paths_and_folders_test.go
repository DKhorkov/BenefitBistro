package paths_and_folders

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"paths_and_folders/testdata"
)

func TestCreateFolderIfNotExists(test *testing.T) {
	workspace_dir, err := getWorkspaceDir()
	if err != nil {
		test.Errorf("Failed to get GO Workspace Dir. Error: %v\n", err)
	}

	test_folder_path := fmt.Sprintf("%v/%v", workspace_dir, testdata.TestsDir) 
	err = CreateFolderIfNotExists(test_folder_path)
	if err != nil {
		if _, err := os.Stat(test_folder_path); os.IsNotExist(err) {
			test.Errorf("Failed to create folder=%v. Error: %v\n", test_folder_path, err)
		}
	}
}

func TestCreateFileIfNotExists(test *testing.T) {
	workspace_dir, err := getWorkspaceDir()
	if err != nil {
		test.Errorf("Failed to get GO Workspace Dir. Error: %v\n", err)
	}

	test_file_path := fmt.Sprintf("%v/%v", workspace_dir, testdata.TestFile) 
	err = CreateFileIfNotExists(test_file_path)
	if err != nil {
		if _, err := os.Stat(test_file_path); os.IsNotExist(err) {
			test.Errorf("Failed to create file=%v. Error: %v\n", test_file_path, err)
		}
	}

}

func getWorkspaceDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	folders_sep := "/"
	splitted_dirs := strings.Split(cwd, folders_sep)
	workspace_dir := strings.Join(
		splitted_dirs[:len(splitted_dirs) - 3], 
		folders_sep) 

	return workspace_dir, nil
}
