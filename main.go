package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/c3s4rfred/sforceds/configs"
	"github.com/c3s4rfred/sforceds/models"
	"github.com/c3s4rfred/sforceds/oauth"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {

	// Authenticate with Salesforce API
	fmt.Println("*****", "Process initiated at", time.Now().String(), "version:", configs.SF_version, "*****")
	fmt.Println("*****", "Trying to login to SalesForce at", time.Now().String(), "*****")
	loginResp, errR := oauth.Login()
	if errR != nil {
		panic(errR)
	} else {
		fmt.Println("*****", "Login success", "*****")
	}

	// Retrieve event data from Salesforce
	err := initSalesForceProcessing(loginResp.AccessToken)
	if err != nil {
		panic(err)
	}
	fmt.Println("*****", "Process terminated at", time.Now().String())
	time.Sleep(5 * time.Second)
}

// initSalesForceProcessing is a method to do all the log extracting process
func initSalesForceProcessing(AccessToken string) error {
	eventURL := configs.InstanceUrl + configs.QueryEndPoint
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
	if err != nil {
		panic(err)
	}
	events := eventsJSON["records"].([]interface{})

	// Prepare event data for SIEM system, first process each EventLogFile returned by salesforce
	for _, event := range events {
		eventMap := event.(map[string]interface{})
		Id := eventMap["Id"].(string)
		err := procLogsById(Id, AccessToken)
		if err != nil {
			log.Println("Error processing log file with ID =", Id, err)
		}
	}
	return nil
}

// procLogsById is a method to get data from a specific EventLogFile
func procLogsById(Id string, token string) error {

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
	errProc := processLogContent(eventResponseBody)
	if errProc != nil {
		return errProc
	}

	return nil
}

// processLogContent is a method to process the logfile from sales force line by line, convert each line in a json object and send it to UTMStack
func processLogContent(data []byte) error {

	reader := csv.NewReader(bytes.NewBuffer(data))

	header, err := reader.Read() // get first line as header
	if err != nil {
		if err != io.EOF {
			return err
		}
	}
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
			log.Println(jsonError)
		} else {
			// Send local_storage to SIEM
			postErr := PostToSiem(jsonEvent)
			if postErr != nil {
				log.Println(postErr)
			}
		}
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
