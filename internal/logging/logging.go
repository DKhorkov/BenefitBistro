package logging

import (
	"config"
	"log"
	"os"
)


var (
	Log *log.Logger
)


func init() {
	create_log_path_if_not_exists()
	var log_file, err = os.OpenFile(config.LogPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
	   panic(err)
	}

	Log = log.New(log_file, "", log.LstdFlags|log.Lshortfile)
}

func create_log_path_if_not_exists() {
	if _, err := os.Open(config.LogDir); os.IsNotExist(err) {
		os.MkdirAll(config.LogDir, os.ModePerm)
	}
	
	if _, err := os.Stat(config.LogPath); os.IsNotExist(err) {
		_, err = os.Create(config.LogPath)

		if err != nil {
			panic(err)
		}
	}
}


func LogTemplateExecuteError(template_name string, err error) {
	Log.Printf("An error occured during trying to execute \"%s\" template: %v\n", template_name, err)
}