package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
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
	Id             int    `json:"id"`
	SourceItem     int    `json:"sourceItem,omitempty"`
	SourceItemName string `json:"sourceItemName,omitempty"`
	SourceTopic    string `json:"sourceTopic,omitempty"`

	TargetItem     int    `json:"targetItem,omitempty"`
	TargetItemName string `json:"targetItemName,omitempty"`
	TargetTopic    string `json:"targetTopic,omitempty"`
}

type MqttStore map[string]*MqttItem

func (m MqttStore) contains(item *MqttItem) bool {
	key := fmt.Sprintf("%s-%d", item.Name, item.Port)
	_, ok := m[key]
	return ok
}

func (m MqttStore) put(item *MqttItem) error {
	key := fmt.Sprintf("%s-%d", item.Name, item.Port)
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
