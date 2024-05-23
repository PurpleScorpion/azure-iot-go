package models

type IotHubConnectionString struct {
	HostName               string `json:"hostName"`
	IotHubConnectionString string `json:"iotHubConnectionString"`
	IotHubName             string `json:"iotHubName"`
	SharedAccessKey        string `json:"sharedAccessKey"`
	SharedAccessKeyName    string `json:"sharedAccessKeyName"`
}
