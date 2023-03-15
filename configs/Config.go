package configs

import (
	"github.com/c3s4rfred/sforceds/utils"
	"time"
)

const (
	// OAuth server variables
	OAuthDialTimeout = 15 * time.Second
	// OAuth header variables
	GrantType = "password"
	// Query definition of the query used to complete QueryEndPoint value
	Query           = "?q=SELECT%20Id,LastModifiedDate%20FROM%20EventLogFile%20WHERE%20LastModifiedDate%20>%20"
	OrderByForQuery = "%20ORDER%20BY%20Id%20ASC"
	// Salesforce constants for json
	SalesForce_dataType = "sales-force"
	// Actual version of the API
	SF_version = "1.0.0 2023-03-15 16:58:06"
)

// State constants
const (
	StateDone                   = "Done"
	StateNext                   = "Next"
	StateTimeYesterday          = "YESTERDAY"
	StateNextEndPointValueUnset = "Unset"
)

// Sales force login params
var (
	// SalesForce connection parameters
	ClientId      = utils.GetEnvOrElse("clientID", "not set")
	ClientSecret  = utils.GetEnvOrElse("clientSecret", "not set")
	Username      = utils.GetEnvOrElse("username", "not set")
	Password      = utils.GetEnvOrElse("password", "not set")
	SecurityToken = utils.GetEnvOrElse("securityToken", "not set")
	InstanceUrl   = utils.GetEnvOrElse("instanceUrl", "not set")
	// OAuth server variables
	// OAuthService represents the salesforce login service url
	OAuthService = utils.GetEnvOrElse("OAuthService", "https://login.salesforce.com")
	// LoginEndpoint represents the salesforce endpoint that returns the oauth2 token based on credentials
	LoginEndpoint = utils.GetEnvOrElse("LoginEndpoint", "/services/oauth2/token")
	// EventsEndPoint represents the salesforce endpoint of the event log files, to get the data by Id
	EventsEndPoint = utils.GetEnvOrElse("EventsEndPoint", "/services/data/v57.0/sobjects/EventLogFile")
	// QueryEndPoint represents the salesforce endpoint to query and return all the event log files endpoints
	QueryEndPoint = utils.GetEnvOrElse("QueryEndPoint", "/services/data/v57.0/query") + Query
	// SIEM URL -> Logs destination
	SiemURL = utils.GetEnvOrElse("siemURL", "http://correlation:8080/v1/newlog")
)
