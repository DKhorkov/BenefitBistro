package paths_and_folders

import (
	"os"
	"strings"
)


func CreateFolderIfNotExists(folder_name string) error {
	if _, err := os.Open(folder_name); os.IsNotExist(err) {
		err = os.MkdirAll(folder_name, os.ModePerm)

		return err
	}

	return nil
}

func CreateFileIfNotExists(file_path string) error {
	if _, err := os.Stat(file_path); os.IsNotExist(err) {
		_, err = os.Create(file_path)

		return err
	}

	return nil
}

func GetWorkspaceDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	folders_sep := "/"
	splitted_dirs := strings.Split(cwd, folders_sep)
	if splitted_dirs[len(splitted_dirs) - 1] == "BenefitBistro" {
		return cwd, nil
	}


	workspace_dir := strings.Join(
		splitted_dirs[:len(splitted_dirs) - 3], 
		folders_sep) 

	return workspace_dir, nil
}

func DeletePath(path string) error {
	err := os.RemoveAll(path) 
    return err
}