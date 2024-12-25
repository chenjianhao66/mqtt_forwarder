package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var mqttStore = MqttStore{}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	engine := gin.Default()

	mqttGroup := engine.Group("/mqtt/client")
	{
		mqttGroup.POST("/add", AddMqttClient)
		mqttGroup.DELETE("/delete", DeleteMqttClient)
		mqttGroup.GET("/list", ListMqttClient)
	}

	forwarderGroup := engine.Group("/mqtt/forwarder")
	{
		forwarderGroup.POST("/add", AddForwarder)
		forwarderGroup.DELETE("/delete", DeleteForwarder)
		forwarderGroup.GET("/list", ListForwarder)
	}
	engine.Run(":8888")
}

func AddMqttClient(ctx *gin.Context) {
	item := new(MqttItem)
	if err := ctx.ShouldBindJSON(item); err != nil {
		log.Println("[添加MQTT客户端] err: ", err)
		ctx.String(http.StatusInternalServerError, "添加错误")
		return
	}
	log.Printf("[添加Mqtt客户端] addr: %s, port: %d, name: %s\n", item.Addr, item.Port, item.Name)
	if mqttStore.contains(item) {
		ctx.String(http.StatusInternalServerError, "该mqtt客户端已经存在，唯一表示是地址和端口的组合")
		return
	}

	if err := mqttStore.put(item); err != nil {
		ctx.String(http.StatusInternalServerError, "添加失败")
		return
	}
	ctx.String(http.StatusOK, "添加成功")
}

func DeleteMqttClient(ctx *gin.Context) {
	req := new(DeleteMqttItemReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.Println("[删除MQTT客户端] 绑定结构体失败，err: ", err)
		ctx.String(http.StatusInternalServerError, "删除错误,绑定结构体失败.")
		return
	}

	log.Printf("[删除Mqtt客户端] addr: %s, port: %d\n", req.Addr, req.Port)
	mqttStore.del(req.Addr, req.Port)
	ctx.String(http.StatusOK, "删除成功")
}

func ListMqttClient(ctx *gin.Context) {
	mqttItems := mqttStore.values()
	ctx.JSON(http.StatusOK, mqttItems)
}

func AddForwarder(ctx *gin.Context) {

}

func DeleteForwarder(ctx *gin.Context) {

}

func ListForwarder(ctx *gin.Context) {

}
