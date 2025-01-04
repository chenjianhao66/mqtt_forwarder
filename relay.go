package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net"
	"time"
)

const (
	DO    = 0x01
	DI    = 0x02
	SetDO = 0x05
)

type ConnectParams struct {
	Name string `json:"name"`
	Addr string `json:"addr,omitempty"`
	Port int    `json:"port,omitempty"`

	currentIndex int
	conn         net.Conn
	ctx          context.Context
	cancelFunc   context.CancelFunc

	DO1 Point `json:"DO1"`
	DO2 Point `json:"DO2"`
	DO3 Point `json:"DO3"`
	DO4 Point `json:"DO4"`
	DO5 Point `json:"DO5"`
	DO6 Point `json:"DO6"`
	DO7 Point `json:"DO7"`
	DO8 Point `json:"DO8"`

	DI1 Point `json:"DI1"`
	DI2 Point `json:"DI2"`
	DI3 Point `json:"DI3"`
	DI4 Point `json:"DI4"`
	DI5 Point `json:"DI5"`
	DI6 Point `json:"DI6"`
	DI7 Point `json:"DI7"`
	DI8 Point `json:"DI8"`
}

func (c *ConnectParams) DIPointStatus() string {
	return fmt.Sprintf("[ DI1:%t,DI2:%t,DI3:%t,DI4:%t,DI5:%t,DI6:%t,DI7:%t,DI8:%t ]", c.DI1.Status, c.DI2.Status, c.DI3.Status, c.DI4.Status, c.DI5.Status, c.DI6.Status, c.DI7.Status, c.DI8.Status)
}

func (c *ConnectParams) DOPointStatus() string {
	return fmt.Sprintf("[ DO1:%t,DO2:%t,DO3:%t,DO4:%t,DO5:%t,DO6:%t,DO7:%t,DO8:%t ]", c.DO1.Status, c.DO2.Status, c.DO3.Status, c.DO4.Status, c.DO5.Status, c.DO6.Status, c.DO7.Status, c.DO8.Status)
}

type Point struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

func (c *ConnectParams) Key() string {
	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}

type Command struct {
	Addr        string `json:"addr,omitempty"`
	Port        int    `json:"port,omitempty"`
	PointNumber int    `json:"pointNumber,omitempty"`
	Status      bool   `json:"status,omitempty"`
}

type RelayStore map[string]*ConnectParams

func (r RelayStore) contains(key string) bool {
	_, ok := r[key]
	return ok
}

var relayStore = RelayStore{}

func Connect(ctx *gin.Context) {
	c := new(ConnectParams)
	if err := ctx.ShouldBindJSON(c); err != nil {
		log.Println("[连接聚英错误] err: ", err)
		FailWithMsg(ctx, "参数错误")
		return
	}

	if relayStore.contains(c.Key()) {
		log.Println("连接聚英错误，ip和port不能重复存在.")
		FailWithMsg(ctx, "连接聚英错误，ip和port不能重复存在.")
		return
	}

	if err := initRelay(c); err != nil {
		FailWithMsg(ctx, "连接聚英失败")
		return
	}
	Success(ctx)
}

func ListRelay(ctx *gin.Context) {
	var result []*ConnectParams
	for _, params := range relayStore {
		result = append(result, params)
	}

	SuccessWithData(ctx, result)
}

func SwitchRelayPointStatus(ctx *gin.Context) {
	c := new(Command)
	if err := ctx.ShouldBindJSON(c); err != nil {
		log.Println("[切换继电器状态] err: ", err)
		FailWithMsg(ctx, "参数错误")
		return
	}

	if !relayStore.contains(fmt.Sprintf("%s:%d", c.Addr, c.Port)) {
		FailWithMsg(ctx, "该聚英不存在")
		return
	}

	params := relayStore[fmt.Sprintf("%s:%d", c.Addr, c.Port)]
	log.Printf("要发送的聚英：%s, DO%d, 要切换的状态：%t", params.Key(), c.PointNumber, c.Status)
	sendSwitchCommand(params, c.PointNumber-1, c.Status)
	SuccessWithData(ctx, params)
}

func Disconnect(ctx *gin.Context) {
	c := new(Command)
	if err := ctx.ShouldBindJSON(c); err != nil {
		log.Println("[断开继电器状态] err: ", err)
		FailWithMsg(ctx, "参数错误")
		return
	}
	key := fmt.Sprintf("%s:%d", c.Addr, c.Port)
	params, ok := relayStore[key]
	if !ok {
		FailWithMsg(ctx, "对应的聚英不存在")
		return
	}
	params.cancelFunc()
	params.conn.Close()
	delete(relayStore, key)
	Success(ctx)
}

func RelayStatusSSE(ctx *gin.Context) {
	command := new(Command)
	if err := ctx.ShouldBindJSON(command); err != nil {
		log.Println("[获取继电器状态失败] err: ", err)
		FailWithMsg(ctx, "参数错误")
		return
	}
	key := fmt.Sprintf("%s:%d", command.Addr, command.Port)
	if !relayStore.contains(key) {
		FailWithMsg(ctx, "对应的聚英不存在")
		return
	}

	// 设置响应头
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	c := ctx.Request.Context()
	ticker := time.NewTicker(time.Second)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-c.Done():
			log.Println("sse 断开连接")
			return
		case <-ticker.C:
			params := relayStore[key]
			ctx.SSEvent("message", params)
			ctx.Writer.Flush()
			log.Println("sse flush.")
		}
	}

}

func initRelay(param *ConnectParams) error {
	dialer := net.Dialer{Timeout: 3 * time.Second}
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", param.Addr, param.Port))
	if err != nil {
		log.Println("连接聚英失败，错误：", err)
		return err
	}
	param.conn = conn
	param.ctx, param.cancelFunc = context.WithCancel(context.Background())
	log.Println("relay 已初始化.", conn.LocalAddr(), conn.RemoteAddr())
	relayStore[param.Key()] = param
	go loopSendQueryCommand(param)
	go receive(param)
	return nil
}

func loopSendQueryCommand(param *ConnectParams) {
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			switch param.currentIndex {
			case 0:
				//log.Printf("向聚英[%s]发送 DO口查询指令", param.Key())
				sendQueryStatusCommand(param, DO)
			case 1:
				//log.Printf("向聚英[%s]发送 DI口查询指令", param.Key())
				sendQueryStatusCommand(param, DI)
			}
			param.currentIndex += 1
			param.currentIndex %= 2
		case <-param.ctx.Done():
			log.Println("销毁聚英")
			return
		}
	}
}

func receive(param *ConnectParams) {

	reader := bufio.NewReader(param.conn)
	for {
		select {
		case <-param.ctx.Done():
			return
		default:
			b := make([]byte, 1024)
			n, err := reader.Read(b)
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Println("读取聚英消息：读取到EOF.")
					break
				}
				log.Println("读取聚英消息失败，err: ", err)
				break
			}
			if n == 0 {
				//log.Printf("读取聚英[%s]消息，长度为0， 跳过.", param.Key())
				continue
			}
			message := b[:n]
			if len(message) < 4 {
				log.Println("读取聚电消息失败，消息长度小于4")
				continue
			}
			command := int(message[1])

			switch command {
			case DO:
				param.DO1.Status = (message[3] & 0x01) != 0
				param.DO2.Status = (message[3] & 0x02) != 0
				param.DO3.Status = (message[3] & 0x04) != 0
				param.DO4.Status = (message[3] & 0x08) != 0
				param.DO5.Status = (message[3] & 0x10) != 0
				param.DO6.Status = (message[3] & 0x20) != 0
				param.DO7.Status = (message[3] & 0x40) != 0
				param.DO8.Status = (message[3] & 0x80) != 0
				//log.Printf("接收到聚英[%s]的DO口查询消息回复: [%d]. \n %s", param.Key(), message[3], param.DOPointStatus())
				//log.Printf("%#x", message)
			case DI:
				param.DI1.Status = (message[3] & 0x01) != 0
				param.DI2.Status = (message[3] & 0x02) != 0
				param.DI3.Status = (message[3] & 0x04) != 0
				param.DI4.Status = (message[3] & 0x08) != 0
				param.DI5.Status = (message[3] & 0x10) != 0
				param.DI6.Status = (message[3] & 0x20) != 0
				param.DI7.Status = (message[3] & 0x40) != 0
				param.DI8.Status = (message[3] & 0x80) != 0
				//log.Printf("接收到聚英[%s]的DI口查询消息回复. \n %s", param.Key(), param.DIPointStatus())
			default:
				log.Printf("收到无效的聚英[%s]回复消息, command: %d", param.Key(), command)
			}
		}
	}

}

func sendQueryStatusCommand(param *ConnectParams, queryType int) {
	command := []byte{0xFE, byte(queryType), 0x00, 0x00, 0x00, 0x08, 0x00, 0x00}
	sendCommand(param, command)
}

func sendSwitchCommand(param *ConnectParams, pointNumber int, opened bool) {
	command := []byte{0xFE, SetDO, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00}
	command[3] = byte(pointNumber & 0xFF)
	if opened {
		command[4] = 0xFF
	} else {
		command[4] = 0x00
	}
	sendCommand(param, command)
}

func sendADCommand(param *ConnectParams) {
	log.Printf("向聚英[%s] 发送模拟量AD查询指令.", param.Key())
	command := []byte{0xFE, 0x04, 0x03, 0xEE, 0x00, 0x08, 0x85, 0xB2}
	sendCommand(param, command)
}

func sendCommand(param *ConnectParams, command []byte) {
	crc16 := getCRC16(command)
	command[len(command)-2] = byte(crc16 & 0xFF)
	command[len(command)-1] = byte((crc16 >> 8) & 0xFF)
	n, err := param.conn.Write(command)
	if err != nil {
		log.Printf("发送聚英[%s] 消息失败，err: %v", param.Key(), err)
		return
	}
	if n != len(command) {
		log.Printf("发送聚英[%s] 消息失败，发送长度和消息长度不一致", param.Key())
		return
	}
	log.Printf("发送聚英[%s] 消息成功", param.Key())
}

func getCRC16(command []byte) int {
	start := 0
	size := len(command) - 2
	var CRC16 = 65535

	for i := 0; i < size; i++ {
		CRC16 ^= int(command[start+i] & 0xFF)
		CRC16 &= 0xFFFF

		for j := 0; j < 8; j++ {
			tmp := CRC16 & 0x0001
			CRC16 >>= 1
			if tmp == 1 {
				CRC16 ^= 0xA001
				CRC16 &= 0xFFFF
			}
		}
	}

	return CRC16 & 0xFFFF
}
