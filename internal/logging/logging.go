package logging

import (
	"log"
	"os"

	"config"
	"paths_and_folders"
)


var (
	Log *log.Logger
)


func init() {
	err := paths_and_folders.CreateFolderIfNotExists(config.LogDir)
	if err != nil {
		panic(err)
	}

	err = paths_and_folders.CreateFileIfNotExists(config.LogPath)
	if err != nil {
		panic(err)
	}

	var log_file *os.File
	log_file, err = os.OpenFile(config.LogPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
	   panic(err)
	}

	Log = log.New(log_file, "", log.LstdFlags|log.Lshortfile)
}

func LogTemplateExecuteError(template_name string, err error) {
	Log.Printf("An error occured during trying to execute \"%s\" template: %v\n", template_name, err)
}
