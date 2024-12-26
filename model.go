package main

import (
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// MqttItem mqtt服务结构体
type MqttItem struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Addr       string `json:"addr"`
	Port       int    `json:"port"`
	NeedVerify bool   `json:"need_verify"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Enable     bool   `json:"enable"`

	client mqtt.Client
}

type DeleteMqttItemReq struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

// MqttForwarderItem mqtt转发规则
type MqttForwarderItem struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Enable bool   `json:"enable"`

	SourceItemAddr string `json:"sourceItemAddr,omitempty"`
	SourceItemPort int    `json:"sourceItemPort,omitempty"`
	SourceItemName string `json:"sourceItemName,omitempty"`
	SourceTopic    string `json:"sourceTopic,omitempty"`

	TargetItemAddr string `json:"targetItemAddr,omitempty"`
	TargetItemPort int    `json:"targetItemPort,omitempty"`
	TargetItem     string `json:"targetItem,omitempty"`
	TargetItemName string `json:"targetItemName,omitempty"`
	TargetTopic    string `json:"targetTopic,omitempty"`
}

type MqttStore map[string]*MqttItem

func (m MqttStore) contains(addr string, port int) bool {
	key := fmt.Sprintf("%s-%d", addr, port)
	_, ok := m[key]
	return ok
}

func (m MqttStore) put(item *MqttItem) error {
	key := fmt.Sprintf("%s-%d", item.Addr, item.Port)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", item.Addr, item.Port))
	opts.SetClientID(item.Name)
	if item.NeedVerify {
		opts.SetUsername(item.Username)
		opts.SetPassword(item.Password)
	}
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(5 * time.Minute)
	opts.SetConnectRetry(true)

	opts.OnConnect = func(client mqtt.Client) {
		log.Printf("[MQTT客户端] %s, %s 已连接", key, item.Name)
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("[MQTT客户端] %s, %s 连接丢失", key, item.Name)
		item.Enable = false
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("[MQTT客户端] 连接失败", token.Error())
		return token.Error()
	}

	item.client = client
	item.Enable = true
	m[key] = item
	return nil
}

func (m MqttStore) get(key string) (*MqttItem, bool) {
	item, ok := m[key]
	return item, ok
}

func (m MqttStore) del(addr string, port int) {
	key := fmt.Sprintf("%s-%d", addr, port)
	if item, ok := m.get(key); ok {
		item.client.Disconnect(250)
		delete(m, key)
	}
}

func (m MqttStore) values() []*MqttItem {
	var items []*MqttItem
	for _, item := range m {
		items = append(items, item)
	}
	return items
}

type ForwarderStore map[int]*MqttForwarderItem

func (f ForwarderStore) add(item *MqttForwarderItem) error {
	source, _ := mqttStore.get(fmt.Sprintf("%s-%d", item.SourceItemAddr, item.SourceItemPort))
	target, _ := mqttStore.get(fmt.Sprintf("%s-%d", item.TargetItemAddr, item.TargetItemPort))

	msgHandler := func(client mqtt.Client, msg mqtt.Message) {
		if token := target.client.Publish(item.TargetTopic, 0, false, msg.Payload()); token.Error() != nil {
			// todo 通过sse发送过前端
			log.Println("[MQTT转发器] 转发消息失败", token.Error())
		}
	}

	if token := source.client.Subscribe(item.SourceTopic, 0, msgHandler); token.Error() != nil {
		log.Println("[MQTT转发器] 订阅topic失败,", token.Error())
		return token.Error()
	}

	id := rand.Intn(1000)

	for _, ok := f[id]; ok; {
		id = rand.Intn(1000)
	}
	item.Id = id
	item.Enable = true
	f[id] = item
	return nil
}

func (f ForwarderStore) del(item *MqttForwarderItem) error {
	source, _ := mqttStore.get(fmt.Sprintf("%s-%d", item.SourceItemAddr, item.SourceItemPort))

	if token := source.client.Unsubscribe(item.SourceTopic); token.Error() != nil {
		log.Println("[MQTT转发器] 取消订阅topic失败,", token.Error())
		return token.Error()
	}

	delete(f, item.Id)
	return nil
}

func (f ForwarderStore) values() []*MqttForwarderItem {
	var items []*MqttForwarderItem
	for _, item := range f {
		items = append(items, item)
	}
	return items
}

func (f ForwarderStore) switchStatus(id int, status bool) error {
	item, ok := f[id]
	if !ok {
		return errors.New("对应的转发器不存在")
	}
	source, ok := mqttStore.get(fmt.Sprintf("%s-%d", item.SourceItemAddr, item.SourceItemPort))
	if !ok {
		return errors.New("对应的mqtt客户端不存在")
	}

	if status {
		target, _ := mqttStore.get(fmt.Sprintf("%s-%d", item.TargetItemAddr, item.TargetItemPort))
		msgHandler := func(client mqtt.Client, msg mqtt.Message) {
			if token := target.client.Publish(item.TargetTopic, 0, false, msg.Payload()); token.Wait() && token.Error() != nil {
				// todo 通过sse发送过前端
				log.Println("[MQTT转发器] 转发消息失败", token.Error())
			}
		}
		if token := source.client.Subscribe(item.SourceTopic, 0, msgHandler); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	} else {
		if token := source.client.Unsubscribe(item.SourceTopic); token.Wait() && token.Error() != nil {
			log.Printf("[关闭MQTT转发器] 名称：%s，关闭失败，错误：%v", item.Name, token.Error())
			return token.Error()
		}
	}
	item.Enable = status
	return nil
}

type Response struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

func Success(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

func SuccessWithData(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": data,
	})
}

func Fail(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusInternalServerError,
		"msg":  "fail",
	})
}

func FailWithMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusInternalServerError,
		"msg":  msg,
	})
}
