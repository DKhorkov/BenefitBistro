package db_adapter

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"config"
	"db_models"
	"logging"
	"paths_and_folders"
)


type DatabaseAdapter struct {
	db *gorm.DB
}

func (db_adapter *DatabaseAdapter) openConnection() error {
	err := paths_and_folders.CreateFolderIfNotExists(config.DatabaseFolder)
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

func (db_adapter *DatabaseAdapter) closeConnection() error {
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
	// TODO подумать над делегатом для открытия соединения
	if err := db_adapter.openConnection(); err != nil {
		logging.Log.Printf("Failed to save Employee via problems with database connecton: %v\n", err)
		return err
	}

	defer db_adapter.closeConnection()

	if _, err := db_adapter.findUserByUsername(username); err == nil {
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
	if err := db_adapter.openConnection(); err != nil {
		logging.Log.Printf("Failed to save Token via problems with database connecton: %v\n", err)
		return err
	}

	defer db_adapter.closeConnection()

	user , err := db_adapter.findUserByUsername(username)
	if err != nil {
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

func (db_adapter *DatabaseAdapter) CompareEmployeeAuthData(username, password string) (bool, error) {
	if err := db_adapter.openConnection(); err != nil {
		logging.Log.Printf("Failed to compare EmployeeAuthData via problems with database connecton: %v\n", err)
		return false, err
	}

	defer db_adapter.closeConnection()

	user, err := db_adapter.findUserByUsername(username)
	if err != nil {
		return false, err
	}

	if user.Password != password {
		logging.Log.Printf("Incorrect password \"%v\" for username \"%v\"\n", password, username)
		return false, errors.New("incorrect password")
	}

	return true, nil
}

func (db_adapter *DatabaseAdapter) findUserByUsername(username string) (*db_models.EmployeeUsers, error) {
	user := &db_models.EmployeeUsers{Username: username}
	if err := db_adapter.db.Where(user).First(user).Error; err != nil {
		logging.Log.Printf("Failed to find user with username=%v. Error:%v\n", username, err)
		return user, err
	}

	return user, nil
}

func (db_adapter *DatabaseAdapter) ValidateEmployeeToken(token string) (*db_models.EmployeeUsers, error) {
	user := &db_models.EmployeeUsers{}
	token_to_find := &db_models.EmployeeTokens{Token: token}

	if err := db_adapter.openConnection(); err != nil {
		logging.Log.Printf("Failed to compare EmployeeAuthData via problems with database connecton: %v\n", err)
		return user, err
	}

	defer db_adapter.closeConnection()

	if err := db_adapter.db.Where(token_to_find).First(token_to_find).Error; err != nil {
		logging.Log.Printf("Failed to find token=%v. Error:%v\n", token, err)
		return user, err
	}
	
	if time.Now().After(token_to_find.Expires) {
		logging.Log.Printf("Token=%v has Expired\n", token)
		db_adapter.deleteEmployeeToken(token)
		return user, errors.New("token expired")
	}

	employee_users_table_name, err := db_adapter.getTableName(&db_models.EmployeeUsers{})
	if err != nil {
		return user, err
	}

	employee_tokens_table_name, err := db_adapter.getTableName(&db_models.EmployeeTokens{})
	if err != nil {
		return user, err
	}

	join_stmt := fmt.Sprintf(
		"JOIN %v ON %v.user_id = %v.id WHERE %v.id = %v", 
		employee_tokens_table_name, 
		employee_tokens_table_name, 
		employee_users_table_name,
		employee_users_table_name,
		token_to_find.UserID,
	)

	if err := db_adapter.db.Table(employee_users_table_name).Select("*").Joins(join_stmt).Scan(user).Error; err != nil {
		logging.Log.Printf("Failed to find user for token=%v. Error:%v\n", token, err)
		return user, err
	}

	return user, nil
}

func (db_adapter *DatabaseAdapter) ValidateHirerToken(token string) (*db_models.HirerUsers, error) {
	return &db_models.HirerUsers{}, nil
}

func (db_adapter *DatabaseAdapter) getTableName(model interface{}) (string, error) {
	stmt := &gorm.Statement{DB: db_adapter.db}
	if err := stmt.Parse(model); err != nil {
		logging.Log.Printf("Failed to find Table name for model=%v\n", model)
		return "", err
	}

	tableName := stmt.Schema.Table
	return tableName, nil
}

func (db_adapter *DatabaseAdapter) DeleteToken(token string) error {
	if strings.HasPrefix(token, config.Token.EmployeePrefix) {
		return db_adapter.deleteEmployeeToken(token)
	} else {
		return db_adapter.deleteHirerToken(token)
	}
}

func (db_adapter *DatabaseAdapter) deleteEmployeeToken(token string) error {
	if err := db_adapter.openConnection(); err != nil {
		logging.Log.Printf("Failed to compare EmployeeAuthData via problems with database connecton: %v\n", err)
		return err
	}

	defer db_adapter.closeConnection()

	token_to_delete := &db_models.EmployeeTokens{Token: token}
	if err := db_adapter.db.Where(token_to_delete).Delete(token_to_delete).Error; err != nil {
		logging.Log.Printf("Failed to delete token=%v. Error: %v\n", token, err)
		return err
	}

	return nil
}

func (db_adapter *DatabaseAdapter) deleteHirerToken(token string) error {
	if err := db_adapter.openConnection(); err != nil {
		logging.Log.Printf("Failed to compare EmployeeAuthData via problems with database connecton: %v\n", err)
		return err
	}

	defer db_adapter.closeConnection()

	return nil
}
