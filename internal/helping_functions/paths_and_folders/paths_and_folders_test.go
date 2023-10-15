package paths_and_folders

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"paths_and_folders/testdata"
)

func TestGetWorkspaceDir(test *testing.T) {
	workspace_dir, err := GetWorkspaceDir()
	assert.True(test, err == nil)
	assert.Equal(
		test, 
		workspace_dir, 
		testdata.ExpectedWorkspaceDir, 
		"TestgetWorkspaceDir failed: %v != %v\n", workspace_dir, testdata.ExpectedWorkspaceDir)
}

func TestCreateFolderIfNotExists(test *testing.T) {
	workspace_dir, err := GetWorkspaceDir()
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
	workspace_dir, err := GetWorkspaceDir()
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

func TestDeletePath(test *testing.T) {
	workspace_dir, err := GetWorkspaceDir()
	if err != nil {
		test.Errorf("Failed to get GO Workspace Dir. Error: %v\n", err)
	}

	test_file_path := fmt.Sprintf("%v/%v", workspace_dir, testdata.TestFile)
	err = DeletePath(test_file_path)
	assert.True(test, err == nil)

	test_folder_path := fmt.Sprintf("%v/%v", workspace_dir, testdata.TestsDir) 
	err = DeletePath(test_folder_path)
	assert.True(test, err == nil)
}