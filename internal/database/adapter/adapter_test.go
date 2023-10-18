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
	db_adapter := prepareTestDatabase(test)
	defer deleteTestDatabase(test)

	has_employee_users_table := db_adapter.db.Migrator().HasTable(&db_models.EmployeeUsers{})
	assert.False(test, has_employee_users_table)

	err := db_adapter.createTableIfNotExists(&db_models.EmployeeUsers{})
	assert.True(test, err == nil)

	has_employee_users_table = db_adapter.db.Migrator().HasTable(&db_models.EmployeeUsers{})
	assert.True(test, has_employee_users_table)
}

func TestSaveEmployee(test *testing.T) {
	db_adapter := prepareTestDatabase(test)
	defer deleteTestDatabase(test)
	
	err := db_adapter.createTableIfNotExists(&db_models.EmployeeUsers{})
	assert.True(test, err == nil)
	user := &db_models.EmployeeUsers{Username: testdata.Username}
	err = db_adapter.db.Where(user).First(user).Error
	assert.False(test, err == nil)

	err = db_adapter.SaveEmployee(testdata.Username, testdata.Password)
	assert.True(test, err == nil)

	// Can not save user if one with same username already exists
	db_adapter = prepareTestDatabase(test)
	err = db_adapter.SaveEmployee(testdata.Username, testdata.Password)
	assert.False(test, err == nil)

	db_adapter = prepareTestDatabase(test)
	err = db_adapter.db.Where(user).First(user).Error
	assert.True(test, err == nil)

	deleteTestDatabase(test)
}

func TestFindUserByUsername(test *testing.T) {
	db_adapter := prepareTestDatabase(test)
	defer deleteTestDatabase(test)

	// User doesn't exist yet
	_, err := db_adapter.findUserByUsername(testdata.Username)
	assert.False(test, err == nil)

	err = db_adapter.SaveEmployee(testdata.Username, testdata.Password)
	assert.True(test, err == nil)

	db_adapter = prepareTestDatabase(test)
	_, err = db_adapter.findUserByUsername(testdata.Username)
	assert.True(test, err == nil)
}

func TestSaveEmployeeToken(test *testing.T) {
	db_adapter := prepareTestDatabase(test)
	defer deleteTestDatabase(test)

	// Can not save token if user does not exist
	err := db_adapter.SaveEmployeeToken(testdata.Username, testdata.EmployeeToken)
	assert.False(test, err == nil)

	err = db_adapter.SaveEmployee(testdata.Username, testdata.Password)
	assert.True(test, err == nil)

	err = db_adapter.SaveEmployeeToken(testdata.Username, testdata.EmployeeToken)
	assert.True(test, err == nil)
}

func TestCompareEmployeeAuthData(test *testing.T) {
	db_adapter := prepareTestDatabase(test)
	defer deleteTestDatabase(test)

	// Can not compare data if user does not exist
	accodrs, err := db_adapter.CompareEmployeeAuthData(testdata.Username, testdata.Password)
	assert.False(test, err == nil)
	assert.False(test, accodrs == true)

	err = db_adapter.SaveEmployee(testdata.Username, testdata.Password)
	assert.True(test, err == nil)

	accodrs, err = db_adapter.CompareEmployeeAuthData(testdata.Username, testdata.InvalidPassword)
	assert.False(test, err == nil)
	assert.False(test, accodrs == true)
	
	accodrs, err = db_adapter.CompareEmployeeAuthData(testdata.Username, testdata.Password)
	assert.True(test, err == nil)
	assert.True(test, accodrs == true)
}

func TestGetTableName(test *testing.T) {
	db_adapter := prepareTestDatabase(test)
	defer deleteTestDatabase(test)

	table_name, err := db_adapter.getTableName(&db_models.EmployeeUsers{})
	assert.True(test, err == nil)
	assert.True(test, table_name == testdata.TableName)
}

func TestValidateEmployeeToken(test *testing.T) {
	db_adapter := prepareTestDatabase(test)
	defer deleteTestDatabase(test)

	// Can't validate token if it doesn't exist
	user, err := db_adapter.ValidateEmployeeToken(testdata.EmployeeToken)
	assert.False(test, err == nil)
	assert.True(test, user.ID == 0)
	assert.True(test, user.Username == "")

	err = db_adapter.SaveEmployee(testdata.Username, testdata.Password)
	assert.True(test, err == nil)

	err = db_adapter.SaveEmployeeToken(testdata.Username, testdata.EmployeeToken)
	assert.True(test, err == nil)

	user, err = db_adapter.ValidateEmployeeToken(testdata.EmployeeToken)
	assert.True(test, err == nil)
	assert.True(test, user.ID == 1)
	assert.True(test, user.Username == testdata.Username)
}

func TestDeleteEmployeeToken(test *testing.T) {
	db_adapter := prepareTestDatabase(test)
	defer deleteTestDatabase(test)

	err := db_adapter.SaveEmployee(testdata.Username, testdata.Password)
	assert.True(test, err == nil)

	err = db_adapter.SaveEmployeeToken(testdata.Username, testdata.EmployeeToken)
	assert.True(test, err == nil)

	db_adapter = prepareTestDatabase(test)

	created_token := &db_models.EmployeeTokens{}
	db_adapter.db.Where(&db_models.EmployeeTokens{Token: testdata.EmployeeToken}).First(created_token)
	assert.True(test, created_token.ID == 1)
	assert.True(test, created_token.UserID == 1)

	err = db_adapter.deleteEmployeeToken(testdata.EmployeeToken)
	assert.True(test, err == nil)

	deleted_token := &db_models.EmployeeTokens{}
	db_adapter.db.Where(&db_models.EmployeeTokens{Token: testdata.EmployeeToken}).First(deleted_token)
	assert.True(test, deleted_token.ID == 0)
	assert.True(test, deleted_token.UserID == 0)
}

func TestDeleteToken(test *testing.T) {
	db_adapter := prepareTestDatabase(test)
	defer deleteTestDatabase(test)

	err := db_adapter.DeleteToken(testdata.EmployeeToken)
	assert.True(test, err == nil)

	err = db_adapter.DeleteToken(testdata.HirerToken)
	assert.True(test, err == nil)

	err = db_adapter.DeleteToken("some_faik_token")
	assert.False(test, err == nil)
}

func prepareTestDatabase(test *testing.T) DatabaseAdapter {
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

	return db_adapter
}

func deleteTestDatabase(test *testing.T) {
	err := paths_and_folders.DeletePath(testdata.DatabaseName)
	assert.True(test, err == nil)
}
