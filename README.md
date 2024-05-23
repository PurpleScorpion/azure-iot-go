## azure旗下iothub消息发送-go版本
`因为项目为go语言,且需要azure旗下iothub进行消息发送, 但是azure并不支持go的该版本, 所以自己写了一个简易版的,仅用于消息的发送`

`注意: eventhub消息订阅请使用官方的sdk: https://github.com/Azure/azure-event-hubs-go`

### 使用方式

 - 引入包
    - go get github.com/PurpleScorpion/azure-iot-go
 - 开始使用
   - 使用案例请看demo包下的Demo.go
   - 详细介绍
     - 创建待发送的消息体 
       ```text
         js := utils.NewJSONObject()
         js.FluentPut("check", "gateway")
         以上为示例 , 只要是个结构体就行
       ```
     - 获取连接字符串
       ```text
         iotHubConnectionString := "your iot hub connection string"
       ```
     - 获取deviceID
       ```text
         deviceID := "your iot hub device id"
       ```
     - 创建设备客户端
       ```text
         deviceClient := iothub.NewDirectMethodsClient(iotHubConnectionString)
       ```
     - 设置消息体/超时时间等参数
       ```text
         options := iothub.DirectMethodRequestOptionsBuilder(js, 10, 10)
       ```
     - 发送消息
       ```text
         methodName := "your iot hub method name"
         response := iothub.Invoke(deviceID, methodName, options, deviceClient)
       ```