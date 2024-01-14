package main

import (
	"fmt"
	"github.com/lxzan/gws"
	"time"
)

const (
	PingInterval = 5 * time.Second
	PingWait     = 10 * time.Second
)

func main() {
	app, _, err := gws.NewClient(&Handler{}, &gws.ClientOption{
		Addr: "ws://127.0.0.1:6666/connect",
	})
	if err != nil {
		return
	}
	_ = app.WritePing(nil)
	go func() {
		time.Sleep(PingWait)
		_ = app.WriteMessage(gws.OpcodeBinary, []byte("hello gws!"))
	}()
	app.ReadLoop()
}

type Handler struct{}

func (c *Handler) OnOpen(socket *gws.Conn) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
}

func (c *Handler) OnClose(socket *gws.Conn, err error) {}

func (c *Handler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	_ = socket.WritePong(nil)
}

func (c *Handler) OnPong(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	fmt.Println("I have receive a pong message!")
	time.Sleep(PingInterval)
	_ = socket.WritePing(nil)
}

func (c *Handler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	fmt.Printf("%d %s\n", message.Opcode, message.Bytes())
	//socket.WriteMessage(message.Opcode, message.Bytes())
}
