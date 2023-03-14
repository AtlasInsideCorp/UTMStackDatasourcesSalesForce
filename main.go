package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/c3s4rfred/sforceds/configs"
	"github.com/c3s4rfred/sforceds/db"
	"github.com/c3s4rfred/sforceds/models"
	"github.com/c3s4rfred/sforceds/oauth"
	"github.com/c3s4rfred/sforceds/utils"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {

	// Authenticate with Salesforce API
	fmt.Println(time.Now().Format(time.RFC3339Nano), "*****", "Process initiated, version:", configs.SF_version, "*****")
	fmt.Println(time.Now().Format(time.RFC3339Nano), "*****", "Trying to login to SalesForce *****")
	loginResp, errR := oauth.Login()
	if errR != nil {
		panic(errR)
	} else {
		fmt.Println(time.Now().Format(time.RFC3339Nano), "*****", "Login success", "*****")
	}

	// Testing database connection
	fmt.Println(time.Now().Format(time.RFC3339Nano), "*****", "Trying to connect to local database *****")
	dbcon, errdb := db.InitDB()
	if errdb != nil {
		panic(errdb)
	}
	// Retrieve event data from Salesforce
	err := initSalesForceProcessing(loginResp.AccessToken, dbcon)
	if err != nil {
		panic(err)
	}
	closeErr := db.CloseDB(dbcon)
	if closeErr != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano), "*****", closeErr.Error(), "*****")
	}
	fmt.Println(time.Now().Format(time.RFC3339Nano), "*****", "Process terminated *****")
	time.Sleep(5 * time.Second)
}

// initSalesForceProcessing is a method to do all the log extracting process
func initSalesForceProcessing(AccessToken string, dbcon *gorm.DB) error {
	// First we have to check the last saved state
	eventURL := ""
	sfLastState, errState := db.GetState(dbcon)
	if errState != nil {
		panic(errState)
	}
	// Set the eventURL according to last state
	if strings.Compare(sfLastState.State, configs.StateDone) == 0 {
		eventURL = configs.InstanceUrl + configs.QueryEndPoint + sfLastState.LastDate + configs.OrderByForQuery
	} else {
		eventURL = configs.InstanceUrl + sfLastState.NextEndPoint
	}
	fmt.Println(time.Now().Format(time.RFC3339Nano), "*****", "Trying to get logs from:", eventURL, "*****")

	eventRequest, err := http.NewRequest("GET", eventURL, nil)
	if err != nil {
		panic(err)
	}
	eventRequest.Header.Set("Content-Type", "application/json")
	eventRequest.Header.Set("Accept", "application/json")
	eventRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AccessToken))
	eventClient := &http.Client{}
	eventResponse, err := eventClient.Do(eventRequest)
	if err != nil {
		panic(err)
	}
	defer eventResponse.Body.Close()
	eventResponseBody, err := io.ReadAll(eventResponse.Body)
	if err != nil {
		panic(err)
	}
	var eventsJSON map[string]interface{}
	err = json.Unmarshal(eventResponseBody, &eventsJSON)
	// Handling bad url, not found and token expired errors
	if err != nil {
		var authJson []models.SalesForceError
		autherr := json.Unmarshal(eventResponseBody, &authJson)
		if autherr != nil {
			panic("Error parsing the event, caused by a bad url")
		} else {
			panic("Error parsing the event, caused by: " + authJson[0].ErrorCode + ", message: " + authJson[0].Message)
		}
	}
	// Get "totalSize" and "done" fields in the response to avoid unnecessary execution and errors
	totalSize := eventsJSON["totalSize"].(float64)
	jsonDone := eventsJSON["done"].(bool)
	if totalSize == 0 {
		fmt.Println(time.Now().Format(time.RFC3339Nano), "*****", "Nothing to process", "*****")
	} else {
		fmt.Println(time.Now().Format(time.RFC3339Nano), "*****", "Beginning to process", totalSize, "log files. All items returned on this query? ->", jsonDone, "*****")
		events := eventsJSON["records"].([]interface{})

		// If the query results are complete, set the state to run TODAY's data in the next iteration, else
		// set the state with the next endpoint of the principal query
		state := db.SfState{}
		if jsonDone {
			state = db.SfState{
				State:        configs.StateDone,
				NextEndPoint: configs.StateNextEndPointValueUnset,
			}
		} else {
			nextEndPoint := eventsJSON["nextRecordsUrl"].(string)
			state = db.SfState{
				State:        configs.StateNext,
				NextEndPoint: nextEndPoint,
			}
		}

		// Prepare event data for SIEM system, first process each EventLogFile returned by salesforce
		for _, event := range events {
			eventMap := event.(map[string]interface{})
			Id := eventMap["Id"].(string)
			LastModifiedDate := eventMap["LastModifiedDate"].(string)
			// Check if Id was processed before to avoid duplicates
			wasProcessed := db.FindByID(dbcon, Id)
			if !wasProcessed {
				errLog := procLogsById(Id, LastModifiedDate, AccessToken, dbcon)
				if errLog != nil {
					fmt.Println(time.Now().Format(time.RFC3339Nano), "*****", "Error processing log file with ID =", Id, errLog, "*****")
				}
			} else {
				fmt.Println(time.Now().Format(time.RFC3339Nano), "*****", "Discarding processed log file with ID =", Id, "*****")
			}

		}
		// Finally, update state
		updateErr := db.UpdateState(dbcon, state)
		if updateErr != nil {
			fmt.Println(time.Now().Format(time.RFC3339Nano), "*****", "Error updating state ->", updateErr, "*****")
		} else {
			fmt.Println(time.Now().Format(time.RFC3339Nano), "*****", "Updating final state to: (State ->", state.State, "), (Next EndPoint ->", state.NextEndPoint, ") *****")
		}
	}

	return nil
}

// procLogsById is a method to get data from a specific EventLogFile
func procLogsById(Id string, LastModifiedDate string, token string, dbcon *gorm.DB) error {

	// Retrieve event data from Salesforce
	eventURL := configs.InstanceUrl + configs.EventsEndPoint + "/" + Id + "/LogFile"
	eventRequest, err := http.NewRequest("GET", eventURL, nil)
	if err != nil {
		return err
	}
	eventRequest.Header.Set("Content-Type", "text/csv")
	eventRequest.Header.Set("Accept", "application/json")
	eventRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	eventClient := &http.Client{}
	eventResponse, err := eventClient.Do(eventRequest)
	if err != nil {
		return err
	}
	defer eventResponse.Body.Close()
	eventResponseBody, err := io.ReadAll(eventResponse.Body)
	if err != nil {
		return err
	}
	errProc := processLogContent(eventResponseBody, Id, LastModifiedDate, dbcon)
	if errProc != nil {
		return errProc
	}

	return nil
}

// processLogContent is a method to process the logfile from sales force line by line, convert each line in a json object and send it to UTMStack
func processLogContent(data []byte, Id string, LastModifiedDate string, dbcon *gorm.DB) error {

	reader := csv.NewReader(bytes.NewBuffer(data))

	header, err := reader.Read() // get first line as header
	if err != nil {
		if err != io.EOF {
			return err
		}
	}
	anyInsertion := false // To control if any row was posted successfully
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		// Convert csv event line to json event
		jsonEvent, jsonError := CreateEvent(header, line)
		if jsonError != nil {
			fmt.Println(jsonError)
		} else {
			// Send data to SIEM
			postErr := PostToSiem(jsonEvent)
			if postErr != nil {
				fmt.Println(postErr)
			} else {
				// Row was posted
				anyInsertion = true
			}
		}
	}
	// Insert the processed Id into the DB if almost one row was posted successfully to SIEM
	if anyInsertion {
		LastMdate, dateErr := time.Parse(time.RFC3339Nano, utils.PrepareDate(LastModifiedDate))
		if dateErr != nil {
			fmt.Println(time.Now().Format(time.RFC3339Nano), "Error parsing date while saving log record -> ", dateErr.Error())
		}
		sflogrec := db.SfLogRecord{
			LogDate: LastMdate,
			Log_id:  Id,
		}
		wasInserted := db.InsertId(dbcon, sflogrec)
		if wasInserted {
			fmt.Println(time.Now().Format(time.RFC3339Nano), "EventLogFile -> ", Id, "processed")
			// Update last state time, in sf_state table: Every time that we process a log file, we update
			// the last date to the LastModifiedDate value of the log file
			sfLogDate := db.SfState{
				LastDate: LastModifiedDate,
			}
			updateErr := db.UpdateState(dbcon, sfLogDate)
			if updateErr != nil {
				fmt.Println(time.Now().Format(time.RFC3339Nano), "Error updating last date in state ->", updateErr)
			}
		}
	} else {
		fmt.Println(time.Now().Format(time.RFC3339Nano), "EventLogFile -> ", Id, "couldn't be processed, reason: 0 rows posted to SIEM, may be the correlation server is down!!!")
	}
	return nil
}

// CreateEvent is a method to convert csv event line to json event
func CreateEvent(header []string, lineOfData []string) ([]byte, error) {
	// Creating the array of local_storage
	EventData := make(map[string]string)
	for i := 0; i < len(header); i++ {
		EventData[header[i]] = lineOfData[i]
	}
	// Creating the log record in json format
	creationTime := time.Now()
	toSend := models.EventLog{
		LogTime: creationTime.Format(time.RFC3339Nano),
		LogGlobal: models.Global{
			LogType: "logx",
		},
		LogDataSource: configs.InstanceUrl,
		LogDataType:   configs.SalesForce_dataType,
		LogxData: models.LogElements{
			EventData: EventData,
		},
	}
	response, err := json.MarshalIndent(toSend, "", "  ")
	if err != nil {
		return nil, err
	}

	return response, nil
}

// PostToSiem is a method to post json local_storage to UTMStack
func PostToSiem(jsonData []byte) error {
	// Send event local_storage to SIEM system
	siemURL := configs.SiemURL
	siemRequest, err := http.NewRequest("POST", siemURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	siemRequest.Header.Set("Content-Type", "application/json")
	siemClient := &http.Client{}
	_, err = siemClient.Do(siemRequest)
	if err != nil {
		return err
	}

	return nil
}
