package iothub

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PurpleScorpion/azure-iot-go/iothub/models"
	"github.com/PurpleScorpion/azure-iot-go/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type IotHubCommon struct {
}

func NewDirectMethodsClient(iothubConnectId string) models.DirectMethodsClient {
	if len(iothubConnectId) == 0 {
		panic("Connection string cannot be null or empty")
	}
	var client models.DirectMethodsClient
	client.IotHubConnectionString = getIotHubConnectionString(iothubConnectId)
	client.HostName = client.IotHubConnectionString.HostName
	return client
}

func NewDirectMethodsModuleClient(iothubConnectId, moduleName string) models.DirectMethodsClient {
	if len(iothubConnectId) == 0 {
		panic("Connection string cannot be null or empty")
	}
	if len(moduleName) == 0 {
		panic("ModuleName cannot be null or empty")
	}
	var client models.DirectMethodsClient
	client.IotHubConnectionString = getIotHubConnectionString(iothubConnectId)
	client.HostName = client.IotHubConnectionString.HostName
	client.ModuleName = moduleName
	return client
}

func Invoke(iothubDeviceId string, methodName string, options models.DirectMethodRequestOptions, client models.DirectMethodsClient) models.DirectMethodResponse {
	if len(iothubDeviceId) == 0 {
		panic("deviceId is empty or null.")
	}
	if options.Payload == nil {
		panic("Payload is empty or null.")
	}

	urlStr := getUrlMethod(client, iothubDeviceId)
	return invokeMethod(urlStr, methodName, options, client)
}

func invokeMethod(url string, methodName string, options models.DirectMethodRequestOptions, client models.DirectMethodsClient) models.DirectMethodResponse {
	jsonStr := toJson(MethodParserBuilder(methodName, options))
	token := getAuthenticationToken(client.IotHubConnectionString)

	httpRequest := httpRequestBuilder(url, jsonStr, token)
	return sendPostRequest(httpRequest)
}

func sendPostRequest(httpRequest models.HttpRequest) models.DirectMethodResponse {
	req, err := http.NewRequest("POST", httpRequest.Url, bytes.NewBuffer(httpRequest.Body))
	if err != nil {
		return errorDirectMethodResponse(err.Error())
	}
	//map循环获取key和value
	for key, value := range httpRequest.Headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errorDirectMethodResponse(err.Error())
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("IotHub Error: " + string(result))
		return errorDirectMethodResponse(string(result))
	}
	js := utils.NewJSONObject()
	js.ParseObject(string(result))
	if strings.Contains(string(result), "status") {
		status, _ := js.GetFloat64("status")
		payload := js.GetJSONObject("payload")
		return models.DirectMethodResponse{
			Status:  int32(status),
			Payload: payload.ToJsonString(),
		}
	} else {
		message, _ := js.GetString("Message")
		jo := utils.NewJSONObject()
		jo.ParseObject(message)
		errorCode, _ := jo.GetFloat64("errorCode")
		msg, _ := jo.GetString("message")
		return models.DirectMethodResponse{
			Status:  int32(errorCode),
			Payload: msg,
		}
	}
}

func errorDirectMethodResponse(err string) models.DirectMethodResponse {
	return models.DirectMethodResponse{
		Status:  500,
		Payload: err,
	}
}

func httpRequestBuilder(url string, body string, authorizationToken string) models.HttpRequest {
	var httpRequest models.HttpRequest
	httpRequest.Body = []byte(body)
	httpRequest.Method = "POST"
	httpRequest.Url = url
	httpRequest.Headers = make(map[string]string)
	httpRequest.Headers["ACCEPT"] = models.ACCEPT_VALUE
	httpRequest.Headers["Content-Type"] = models.CONTENT_TYPE
	httpRequest.Headers["Authorization"] = authorizationToken
	httpRequest.Headers["Request-Id"] = "1"
	httpRequest.Headers["User-Agent"] = getUserAgent()
	httpRequest.Headers["Content-Length"] = fmt.Sprintf("%d", len(httpRequest.Body))
	// js, _ := json.Marshal(httpRequest)
	return httpRequest
}

func getUserAgent() string {
	// artifactId := "iot-service-client"
	// groupId := "com.microsoft.azure.sdk.iot"
	version := "2.1.6"
	JAVA_RUNTIME := "1.8.0_152"
	OPERATING_SYSTEM := os.Getenv("OS")
	PROCESSOR_ARCHITECTURE := os.Getenv("GOARCH")

	USER_AGENT_STRING := "com.microsoft.azure.sdk.iot.iot-service-client/" + version + " (" + JAVA_RUNTIME + "; " + OPERATING_SYSTEM + "; " + PROCESSOR_ARCHITECTURE + ")"

	return USER_AGENT_STRING
}

func MethodParserBuilder(methodName string, options models.DirectMethodRequestOptions) models.MethodParser {
	var methodParser models.MethodParser
	methodParser.Name = methodName
	methodParser.ConnectTimeout = options.MethodConnectTimeoutSeconds
	methodParser.ResponseTimeout = options.MethodResponseTimeoutSeconds
	methodParser.Payload = options.Payload
	methodParser.Operation = "invoke"
	return methodParser
}

func toJson(methodParser models.MethodParser) string {
	JsonObject := make(map[string]interface{})
	data, _ := json.Marshal(methodParser.Payload)
	if methodParser.Operation == "invoke" {
		JsonObject["methodName"] = methodParser.Name
		JsonObject["responseTimeoutInSeconds"] = methodParser.ResponseTimeout
		JsonObject["connectTimeoutInSeconds"] = methodParser.ConnectTimeout
		JsonObject["payload"] = string(data)
	} else if methodParser.Operation == "response" {
		panic("response operation not implemented yet.")
	} else if methodParser.Operation == "payload" {
		panic("payload operation not implemented yet.")
	}
	jsonString, _ := json.Marshal(JsonObject)
	return string(jsonString)
}

func getUrlMethod(client models.DirectMethodsClient, deviceId string) string {
	if len(client.HostName) == 0 {
		panic("hostName is empty or null.")
	}
	if len(deviceId) == 0 {
		panic("deviceId is empty or null.")
	}
	if len(client.ModuleName) == 0 {
		return "https://" + client.HostName + "/" + "twins" + "/" + deviceId + "/methods" + "?" + "api-version=2021-04-12"
	}
	//https://<iothubName>.azure-devices.net/twins/<deviceId>/modules/<moduleName>/methods?api-version=2021-04-12
	return "https://" + client.HostName + "/" + "twins" + "/" + deviceId + "/modules/" + client.ModuleName + "/methods" + "?" + "api-version=2021-04-12"
}

func getIotHubConnectionString(iothubConnectId string) models.IotHubConnectionString {
	if len(iothubConnectId) == 0 {
		panic("Connection string cannot be null or empty")
	}
	var connectionString models.IotHubConnectionString
	connectionString.IotHubConnectionString = iothubConnectId
	array := strings.Split(iothubConnectId, ";")
	for _, v := range array {
		if strings.Contains(v, "HostName=") {
			connectionString.HostName = strings.Replace(v, "HostName=", "", -1)
		} else if strings.Contains(v, "SharedAccessKeyName=") {
			connectionString.SharedAccessKeyName = strings.Replace(v, "SharedAccessKeyName=", "", -1)
		} else if strings.Contains(v, "SharedAccessKey=") {
			connectionString.SharedAccessKey = strings.Replace(v, "SharedAccessKey=", "", -1)
		} else if strings.Contains(v, "IotHubName=") {
			connectionString.IotHubName = strings.Replace(v, "IotHubName=", "", -1)
		}
	}
	return connectionString
}

func DirectMethodRequestOptionsBuilder(payload interface{}, methodConnectTimeoutSeconds int32, methodResponseTimeoutSeconds int32) models.DirectMethodRequestOptions {
	var options models.DirectMethodRequestOptions
	if methodConnectTimeoutSeconds > 0 {
		options.MethodConnectTimeoutSeconds = methodConnectTimeoutSeconds
	} else {
		options.MethodConnectTimeoutSeconds = 5
	}

	if methodResponseTimeoutSeconds > 5 {
		options.MethodResponseTimeoutSeconds = methodResponseTimeoutSeconds
	} else {
		options.MethodResponseTimeoutSeconds = 5
	}
	options.Payload = payload
	return options
}

func getAuthenticationToken(connectionString models.IotHubConnectionString) string {
	var iotHubServiceSasToken models.IotHubServiceSasToken
	iotHubServiceSasToken.TokenLifespanSeconds = 3600
	iotHubServiceSasToken.ResourceUri = connectionString.HostName
	iotHubServiceSasToken.KeyValue = connectionString.SharedAccessKey
	iotHubServiceSasToken.KeyName = connectionString.SharedAccessKeyName
	iotHubServiceSasToken.ExpiryTimeSeconds = buildExpiresOn(iotHubServiceSasToken.TokenLifespanSeconds)
	token := buildToken(iotHubServiceSasToken)
	return token
}

func buildToken(iotHubServiceSasToken models.IotHubServiceSasToken) string {

	sharedAccessSignature, _ := NewSharedAccessSignature(
		iotHubServiceSasToken.ResourceUri, iotHubServiceSasToken.KeyName, iotHubServiceSasToken.KeyValue, iotHubServiceSasToken.ExpiryTimeSeconds)
	signature := sharedAccessSignature.Sig
	targetUri := iotHubServiceSasToken.ResourceUri
	keyName := iotHubServiceSasToken.KeyName

	encodedStr := url.QueryEscape(signature)

	signatureString := "SharedAccessSignature sr=" + targetUri + "&sig=" + encodedStr + "&se=" + fmt.Sprint(iotHubServiceSasToken.ExpiryTimeSeconds) + "&skn=" + keyName

	return signatureString
}

func NewSharedAccessSignature(
	resource, policy, key string, expiry int64,
) (*SharedAccessSignature, error) {
	sig, err := mksig(resource, key, expiry)
	if err != nil {
		return nil, err
	}
	return &SharedAccessSignature{
		Sr:  resource,
		Sig: sig,
		Se:  expiry,
		Skn: policy,
	}, nil
}

func mksig(sr, key string, se int64) (string, error) {
	b, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	h := hmac.New(sha256.New, b)
	if _, err := fmt.Fprintf(h, "%s\n%d", url.QueryEscape(sr), se); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func buildExpiresOn(tokenLifespanSeconds int64) int64 {
	expiresOnDate := time.Now().Unix()
	expiresOnDate += tokenLifespanSeconds
	return expiresOnDate
}

type SharedAccessSignature struct {
	Sr  string
	Sig string
	Se  int64
	Skn string
}
