package db_adapter

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"

	"config"
	"db_adapter/testdata"
	"db_models"
	"paths_and_folders"
)

func TestGetDatabaseFolder(test *testing.T) {
	db_adapter := DatabaseAdapter{}
	db_folder := db_adapter.getDatabaseFolder()
	assert.Equal(
		test, 
		db_folder, 
		config.DatabaseFolder, 
		"TestGetDatabaseFolder failed: %v != %v\n", db_folder, config.DatabaseFolder)

	db_adapter = DatabaseAdapter{
		DatabaseFolder: testdata.DatabaseFolder,
	}
	db_folder = db_adapter.getDatabaseFolder()
	assert.Equal(
		test, 
		db_folder, 
		db_adapter.DatabaseFolder, 
		"TestGetDatabaseFolder failed: %v != %v\n", db_folder, testdata.DatabaseFolder)
}

func TestGetDatabaseName(test *testing.T) {
	db_adapter := DatabaseAdapter{}
	db_name := db_adapter.getDatabaseName()
	assert.Equal(
		test, 
		db_name, 
		config.DatabaseName, 
		"TestGetDatabaseName failed: %v != %v\n", db_name, config.DatabaseName)

	db_adapter = DatabaseAdapter{
		DatabaseName: testdata.DatabaseName,
	}
	db_name = db_adapter.getDatabaseName()
	assert.Equal(
		test, 
		db_name, 
		db_adapter.DatabaseName, 
		"TestGetDatabaseName failed: %v != %v\n", db_name, testdata.DatabaseName)
}

func TestCreateTableIfNotExists(test *testing.T) {
	db_adapter := DatabaseAdapter{
		DatabaseFolder: testdata.DatabaseFolder,
		DatabaseName: testdata.DatabaseName,
	}
	
	db_folder := db_adapter.getDatabaseFolder()
	err := paths_and_folders.CreateFolderIfNotExists(db_folder)
	assert.True(test, err == nil)


	db_name := db_adapter.getDatabaseName()
	db_adapter.db, err = gorm.Open(sqlite.Open(db_name), &gorm.Config{})
	assert.True(test, err == nil)

	has_employee_users_table := db_adapter.db.Migrator().HasTable(&db_models.EmployeeUsers{})
	assert.False(test, has_employee_users_table)

	err = db_adapter.createTableIfNotExists(&db_models.EmployeeUsers{})
	assert.True(test, err == nil)

	has_employee_users_table = db_adapter.db.Migrator().HasTable(&db_models.EmployeeUsers{})
	assert.True(test, has_employee_users_table)

	err = paths_and_folders.DeletePath(testdata.DatabaseName)
	assert.True(test, err == nil)
}
