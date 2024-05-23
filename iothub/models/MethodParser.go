package models

type MethodParser struct {
	Name            string      `json:"name"`
	ResponseTimeout int32       `json:"responseTimeout"`
	ConnectTimeout  int32       `json:"connectTimeout"`
	Payload         interface{} `json:"payload"`
	Operation       string      `json:"operation"`
}
