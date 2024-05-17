package models

type IotHubServiceSasToken struct {
	ResourceUri          string `json:"resourceUri"`
	KeyValue             string `json:"keyValue"`
	ExpiryTimeSeconds    int64  `json:"expiryTimeSeconds"`
	KeyName              string `json:"keyName"`
	Token                string `json:"token"`
	TokenLifespanSeconds int64  `json:"tokenLifespanSeconds"`
}
