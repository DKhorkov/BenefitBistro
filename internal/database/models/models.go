package db_models

import (
	"time"

	"gorm.io/gorm"
)


type EmployeeUsers struct {
	gorm.Model
	Username, Password string
}

type EmployeeTokens struct {
	gorm.Model
	UserID int
	EmployeeUsers EmployeeUsers `gorm:"foreignKey:UserID;references:ID"`
	Token string
	Expires time.Time
}

type HirerUsers struct {
	gorm.Model
	Username, Password string
}
