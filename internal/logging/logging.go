package logging

import (
	"log"
	"os"

	"config"
	"helping_functions"
)


var (
	Log *log.Logger
)


func init() {
	err := helping_functions.CreateFolderIfNotExists(config.LogDir)
	if err != nil {
		panic(err)
	}

	err = helping_functions.CreateFileIfNotExists(config.LogPath)
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
