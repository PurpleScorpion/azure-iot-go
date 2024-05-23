package demo

import (
	"fmt"
	"github.com/PurpleScorpion/azure-iot-go/iothub"
	"github.com/PurpleScorpion/azure-iot-go/iothub/models"
)

func demo1() {
	device := Device{Id: 1, DeviceCode: "123", DeviceName: "test"}

	iotPojo := BuildIotMsg("deviceUpdate", device, "system")
	InvokeMethod(iotPojo, "your methodName")
}

type Device struct {
	Id         int32  `json:"id"`
	DeviceCode string `json:"deviceCode"`
	DeviceName string `json:"deviceName"`
}

type IotSendVO struct {
	Operation  string      `json:"operation"`
	ClientID   string      `json:"clientID"`
	DeviceType string      `json:"deviceType"`
	Data       interface{} `json:"data"`
}

// 消息体构建器 , 里面的内容和参数根据自己的需求定义
func BuildIotMsg(Operation string, Data interface{}, deviceType string) IotSendVO {
	var pojo IotSendVO
	pojo.ClientID = "CloudEMS"
	pojo.Operation = Operation
	pojo.DeviceType = deviceType
	pojo.Data = Data
	return pojo
}

/*
根执行器 , 调用此函数进行Iothub消息发送

	payload: 消息体, 由BuildIotMsg构建
	methodName: iothub必要的methodName
*/
func InvokeMethod(payload interface{}, methodName string) models.DirectMethodResponse {
	iotHubConnectionString := getConnectId()

	if isEmpty(iotHubConnectionString) {
		panic("iotHubConnectionString is empty")
	}

	deviceID := getDeviceId()
	if isEmpty(deviceID) {
		panic("deviceID is empty")
	}

	// // 创建设备客户端
	deviceClient := iothub.NewDirectMethodsClient(iotHubConnectionString)

	options := iothub.DirectMethodRequestOptionsBuilder(payload, 10, 10)

	// // 发送方法调用请求
	response := iothub.Invoke(deviceID, methodName, options, deviceClient)
	fmt.Println("----------iothub------------", response.Status, response.Payload)
	return response
}

func getConnectId() string {
	iothubConnectIds := "your iothubConnectIds"
	return iothubConnectIds
}

func getDeviceId() string {
	iothubDeviceIds := "your iothubDeviceId"

	return iothubDeviceIds
}

func isEmpty(str string) bool {
	return len(str) == 0
}
