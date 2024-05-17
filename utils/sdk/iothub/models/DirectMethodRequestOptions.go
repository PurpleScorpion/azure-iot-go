package models

type DirectMethodRequestOptions struct {
	Payload                      interface{} `json:"payload"`
	MethodConnectTimeoutSeconds  int32       `json:"methodConnectTimeoutSeconds"`
	MethodResponseTimeoutSeconds int32       `json:"methodResponseTimeoutSeconds"`
}
