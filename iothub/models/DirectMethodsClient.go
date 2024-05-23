package models

type DirectMethodsClient struct {
	HostName               string                     `json:"hostName"`
	RequestId              int64                      `json:"requestId"`
	Options                DirectMethodRequestOptions `json:"options"`
	IotHubConnectionString IotHubConnectionString     `json:"iotHubConnectionString"`
}
