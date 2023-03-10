package db

import (
	"errors"
	"fmt"
	"github.com/c3s4rfred/sforceds/configs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"
)

// SfLogRecord is the table to store processed Ids
type SfLogRecord struct {
	gorm.Model
	Log_id string
}

// SfState is the table to store last states when processing logfiles
type SfState struct {
	gorm.Model
	LastDate     string // Last date stored of the processed event files
	State        string // Last state of the query executed, can be (Done if the query has no more elements) or (Next, if points to next endpoint of data)
	NextEndPoint string // Used to store the URI of the next endpoint data of the principal query
	/*
		Note: On salesforce, query gives up to 2000 rows, if the result is bigger, then return a new field with a new endpoint
		to call to retrieve the next bunch of data
	*/
}

// InitDB is a method to initialize the DB
func InitDB() (*gorm.DB, error) {
	if _, err := os.Stat("local_storage"); err == nil {
		// Exists
	} else if errors.Is(err, os.ErrNotExist) {
		os.Mkdir("local_storage", 0700)
	}
	db, err := gorm.Open(sqlite.Open("local_storage/sf_processed_logs.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.New("Unable to connect to database, check if the path -> local_storage/sf_processed_logs.db, exists and you have (read/write) permissions")
	}

	// Migrate the schema
	db.AutoMigrate(&SfLogRecord{})
	db.AutoMigrate(&SfState{})
	SetInitState(db)
	fmt.Println(time.Now().String(), "*****", "Database connected", "*****")
	return db, nil
}

// FindByID is a method to search by Id and return SfLogRecord
func FindByID(db *gorm.DB, Id string) bool {
	// Read
	var sflog SfLogRecord
	result := db.Limit(1).Find(&sflog, "Log_id = ?", Id)
	if result.RowsAffected > 0 {
		return true
	}
	return false
}

// InsertId is a method to insert an element into the DB
func InsertId(dbcon *gorm.DB, Id string) bool {
	sfrecord := SfLogRecord{Log_id: Id}
	result := dbcon.Create(&sfrecord)
	if result.RowsAffected > 0 {
		return true
	}
	return false
}

// SetInitState is a method to create the initial state of execution
func SetInitState(db *gorm.DB) {
	var sflog SfState
	result := db.Limit(1).Find(&sflog, "id = 1")
	if result.RowsAffected == 0 {
		// Create initial state row
		db.Create(&SfState{
			LastDate:     configs.StateTimeYesterday,
			State:        configs.StateDone,
			NextEndPoint: configs.StateNextEndPointValueUnset,
		})
	}
}

// UpdateState is a method to update state of execution
func UpdateState(db *gorm.DB, state SfState) error {
	var sflog SfState
	result := db.Limit(1).Find(&sflog, "id = 1")
	if result.RowsAffected == 0 {
		return errors.New("Something odd happens, state not found")
	} else {
		db.Model(&sflog).Updates(state)
	}
	return nil
}

// GetState is a method to get last state of execution
func GetState(db *gorm.DB) (SfState, error) {
	var sflog SfState
	result := db.Limit(1).Find(&sflog, "id = 1")
	if result.RowsAffected == 0 {
		return sflog, errors.New("Unable to get last state, state record not found")
	} else {
		return sflog, nil
	}
}
