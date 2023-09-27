package helping_functions

import (
	"os"
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