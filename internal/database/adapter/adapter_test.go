package db_adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"config"
)

func TestGetDatabaseFolder(test *testing.T) {
	db_adapter := DatabaseAdapter{}
	db_folder := db_adapter.getDatabaseFolder()
	assert.Equal(
		test, 
		db_folder, 
		config.DatabaseFolder, 
		"TestGetDatabaseFolder failed: %v != %v", db_folder, config.DatabaseFolder)

	db_adapter = DatabaseAdapter{
		DatabaseFolder: "tmp/tests_folder/database_folder/",
	}
	db_folder = db_adapter.getDatabaseFolder()
	assert.Equal(
		test, 
		db_folder, 
		db_adapter.DatabaseFolder, 
		"TestGetDatabaseFolder failed: %v != %v", db_folder, config.DatabaseFolder)
}

func TestGetDatabaseName(test *testing.T) {
	db_adapter := DatabaseAdapter{}
	db_name := db_adapter.getDatabaseName()
	assert.Equal(
		test, 
		db_name, 
		config.DatabaseName, 
		"TestGetDatabaseName failed: %v != %v", db_name, config.DatabaseName)

	db_adapter = DatabaseAdapter{
		DatabaseName: "tmp/tests_folder/database_folder/my_test_database.db",
	}
	db_name = db_adapter.getDatabaseName()
	assert.Equal(
		test, 
		db_name, 
		db_adapter.DatabaseName, 
		"TestGetDatabaseName failed: %v != %v", db_name, config.DatabaseName)
}
