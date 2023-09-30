package db_adapter

import (
	"reflect"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"config"
	"db_models"
	"helping_functions"
	"logging"
)


type DatabaseAdapter struct {
	db *gorm.DB
}

func (db_adapter *DatabaseAdapter) OpenConnection() error {
	err := helping_functions.CreateFolderIfNotExists(config.DatabaseFolder)
	if err != nil {
		logging.Log.Printf("Failed to MakeAllFolders %s: %v", config.DatabaseFolder, err)
		panic(err)
	}

	db_adapter.db, err = gorm.Open(sqlite.Open(config.DatabaseName), &gorm.Config{})

	if err != nil {
		logging.Log.Printf("Error during trying to connect to database: %v", err)
		return err
	}

	err = db_adapter.createTablesIfNotExists()
	if err != nil {
		return err
	}

	return nil
}

func (db_adapter *DatabaseAdapter) CloseConnection() error {
	dbInstance, err := db_adapter.db.DB()

	if err != nil {
		logging.Log.Printf("Failed to create instanse for closing DB : %v", err)
		return err
	}

	err = dbInstance.Close()

	if err != nil {
		logging.Log.Printf("Failed to close DB : %v", err)
		return err
	}

	return nil
}

func (db_adapter *DatabaseAdapter) createTablesIfNotExists() error {
	// Список всех таблиц, чтобы создать их итеративно:
	tables := []interface{} {
		&db_models.EmployeeUsers{}, 
		&db_models.HirerUsers{},
		&db_models.EmployeeTokens{},
	}

	for _, table := range tables {
		err := db_adapter.createTableIfNotExists(table)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db_adapter *DatabaseAdapter) createTableIfNotExists(table interface{}) error {
	if !db_adapter.db.Migrator().HasTable(table) {
		err := db_adapter.db.Migrator().CreateTable(table)
		if err != nil {
			logging.Log.Printf(
				"Failed to create table \"%v\" : %v", 
				reflect.TypeOf(table).String(),
				err)
			return err
		}
	}

	return nil
}

func (db_adapter *DatabaseAdapter) SaveEmployee(username, password string) error {
	if err := db_adapter.db.Where(
		"username = ?",
		username).First(
			&db_models.EmployeeUsers{}).Error; err == nil {
				logging.Log.Printf("User with username=%s already exists\n", username)
				return gorm.ErrRegistered
			}

	if err := db_adapter.db.Create(
		&db_models.EmployeeUsers{
			Username: username,
			Password: password,
		}).Error; err != nil {
			logging.Log.Printf("Failed to save EmployeeUser : %v", err)
			return err
	}

	return nil
}

func (db_adapter *DatabaseAdapter) SaveToken(username, token string) error {
	user := &db_models.EmployeeUsers{Username: username}
	if err := db_adapter.db.Where(user).First(user).Error; err != nil {
		logging.Log.Printf("Failed to find user with username=%v. Error:%v\n", username, err)
		return err
	}

	if err := db_adapter.db.Create(
		&db_models.EmployeeTokens{
			Token: token,
			Expires: time.Now().Local().Add(config.Token.ExpiresDuration),
			UserID: int(user.ID),
		}).Error; err != nil {
			logging.Log.Printf("Failed to save EmployeeToken : %v", err)
			return err
		}

	return nil
}

func (db_adapter *DatabaseAdapter) CompareAuthData(username, password string) (bool, error) {
	return true, nil
}
