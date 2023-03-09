package models

// Global Type to insert in the json root before send, always is 'logx'
type Global struct {
	LogType string `json:"type"`
}

// LogElements Sub-element representing the field and its values
type LogElements struct {
	EventData map[string]string `json:"sales_force"`
}

// EventLog Structure of the json log to send to UTMStack
type EventLog struct {
	LogTime       string      `json:"@timestamp"`
	LogGlobal     Global      `json:"global"`
	LogDataSource string      `json:"dataSource"`
	LogDataType   string      `json:"dataType"`
	LogxData      LogElements `json:"logx"`
}
