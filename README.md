## azure旗下iothub消息发送-go版本
`因为项目为go语言,且需要azure旗下iothub进行消息发送, 但是azure并不支持go的该版本, 所以自己写了一个简易版的,仅用于消息的发送`

`注意: eventhub消息订阅请使用官方的sdk: https://github.com/Azure/azure-event-hubs-go`

### 使用方式

 - 复制文件
    - 将utils下的所有文件复制到你的项目中
     - 注意需要更改包的引用
     - [IotHubCommon.go] `azure-iot-go/utils` `azure-iot-go/utils/sdk/iothub/models`
     - 将以上两个包改为复制到你项目后的地址,你的IDE应该会自动更改
 - 开始使用
   - 使用案例请看demo包下的Demo.go
   - 里面详细介绍了该sdk的使用方式