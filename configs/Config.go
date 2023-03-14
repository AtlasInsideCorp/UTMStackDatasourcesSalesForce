package configs

import (
	"os"
	"time"
)

const (
	// OAuth server variables
	OAuthService     = "https://login.salesforce.com"
	LoginEndpoint    = "/services/oauth2/token"
	OAuthDialTimeout = 15 * time.Second
	// OAuth header variables
	GrantType = "password"
	// Events endpoint
	EventsEndPoint = "/services/data/v57.0/sobjects/EventLogFile"
	// Query endpoint
	QueryEndPoint   = "/services/data/v57.0/query?q=SELECT%20Id,LastModifiedDate%20FROM%20EventLogFile%20WHERE%20LastModifiedDate%20>%20"
	OrderByForQuery = "%20ORDER%20BY%20Id%20ASC"
	// Salesforce constants for json
	SalesForce_dataType = "sales-force"
	// SIEM URL
	SiemURL = "http://correlation:8080/v1/newlog"
	// Actual version of the API
	SF_version = "1.0.0 2023-03-14 15:08:24"
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
	ClientId      = os.Getenv("clientID")
	ClientSecret  = os.Getenv("clientSecret")
	Username      = os.Getenv("username")
	Password      = os.Getenv("password")
	SecurityToken = os.Getenv("securityToken")
	InstanceUrl   = os.Getenv("instanceUrl")
)
