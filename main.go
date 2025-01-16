package main

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"log"
	"net/http"
)

var (
	mqttStore      = MqttStore{}
	forwarderStore = ForwarderStore{}
)

//go:embed ui/dist/*
var staticFiles embed.FS

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
		forwarderGroup.POST("/enable/:id", EnableForwarder)
		forwarderGroup.POST("/disable/:id", DisableForwarder)
	}

	relayGroup := engine.Group("/mqtt/relay")
	{
		relayGroup.POST("/connect", Connect)
		relayGroup.GET("/list", ListRelay)
		relayGroup.POST("/command", SwitchRelayPointStatus)
		relayGroup.POST("/disconnect", Disconnect)
		relayGroup.GET("/status", RelayStatusSSE)
	}

	engine.Use(static.Serve("/", static.EmbedFolder(staticFiles, "ui/dist")))
	engine.NoRoute(func(context *gin.Context) {
		fmt.Printf("%s doesn't exists, redirect on /\n", context.Request.URL.Path)
		context.Redirect(http.StatusMovedPermanently, "/")
	})
	engine.Run(":8888")
}

func AddMqttClient(ctx *gin.Context) {
	item := new(MqttItem)
	if err := ctx.ShouldBindJSON(item); err != nil {
		log.Println("[添加MQTT客户端] err: ", err)
		Fail(ctx)
		return
	}
	log.Printf("[添加Mqtt客户端] addr: %s, port: %d, name: %s\n", item.Addr, item.Port, item.Name)
	if mqttStore.contains(item.Addr, item.Port) {
		FailWithMsg(ctx, "该mqtt客户端已经存在，唯一表示是地址和端口的组合")
		return
	}

	if err := mqttStore.put(item); err != nil {
		FailWithMsg(ctx, "添加失败")
		return
	}
	Success(ctx)
}

func DeleteMqttClient(ctx *gin.Context) {
	req := new(DeleteMqttItemReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.Println("[删除MQTT客户端] 绑定结构体失败，err: ", err)
		FailWithMsg(ctx, "删除错误,绑定结构体失败.")
		return
	}

	log.Printf("[删除Mqtt客户端] addr: %s, port: %d\n", req.Addr, req.Port)
	mqttStore.del(req.Addr, req.Port)
	Success(ctx)
}

func ListMqttClient(ctx *gin.Context) {
	mqttItems := mqttStore.values()
	SuccessWithData(ctx, mqttItems)
}

func AddForwarder(ctx *gin.Context) {
	req := new(MqttForwarderItem)
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.Println("[添加转发器] 绑定结构体失败，err: ", err)
		FailWithMsg(ctx, "添加转发器错误,绑定结构体失败.")
		return
	}
	if !mqttStore.contains(req.SourceItemAddr, req.SourceItemPort) {
		log.Println("[添加转发器] 源客户端不存在，请先添加源客户端")
		FailWithMsg(ctx, "添加转发器错误,源客户端不存在.")
		return
	}
	if !mqttStore.contains(req.TargetItemAddr, req.TargetItemPort) {
		log.Println("[添加转发器] 目标客户端不存在，请先添加目标客户端")
		FailWithMsg(ctx, "添加转发器错误,目标客户端不存在.")
		return
	}

	if err := forwarderStore.add(req); err != nil {
		log.Println("[添加转发器] 添加转发器失败，err: ", err)
		FailWithMsg(ctx, "添加转发器错误,添加转发器失败.")
		return
	}

	Success(ctx)
}

func DeleteForwarder(ctx *gin.Context) {
	req := new(MqttForwarderItem)
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.Println("[删除转发器] 绑定结构体失败，err: ", err)
		FailWithMsg(ctx, "删除转发器错误,绑定结构体失败.")
		return
	}

	if err := forwarderStore.del(req); err != nil {
		log.Println("[删除转发器] 删除转发器失败，err: ", err)
		FailWithMsg(ctx, "删除转发器错误,删除转发器失败.")
		return
	}
	Success(ctx)
}

func ListForwarder(ctx *gin.Context) {
	values := forwarderStore.values()
	SuccessWithData(ctx, values)
}

func EnableForwarder(ctx *gin.Context) {
	id := cast.ToInt(ctx.Param("id"))

	if err := forwarderStore.switchStatus(id, true); err != nil {
		log.Println(err)
		FailWithMsg(ctx, "启用转发器失败")
	}

	Success(ctx)
}

func DisableForwarder(ctx *gin.Context) {
	id := cast.ToInt(ctx.Param("id"))

	if err := forwarderStore.switchStatus(id, false); err != nil {
		log.Println(err)
		FailWithMsg(ctx, "关闭转发器失败")
	}

	Success(ctx)
}
