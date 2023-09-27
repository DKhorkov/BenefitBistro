package db_models

import (
	"gorm.io/gorm"
)


type EmployeeUsers struct {
	gorm.Model
	Username, Password string
}

type HirerUsers struct {
	gorm.Model
	Username, Password string
}