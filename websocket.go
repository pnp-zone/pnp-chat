package main

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	"io"
)

type WebSocket struct {
	*websocket.Conn
	Received chan string
	Logger   echo.Logger
}

func NewWebSocket(c echo.Context, ws *websocket.Conn) *WebSocket {
	socket := &WebSocket{
		ws,
		make(chan string),
		c.Logger(),
	}

	// Start a goroutine listening for messages and sending them to the Received channel
	go func() {
		defer close(socket.Received)
		for {
			msg := ""
			err := websocket.Message.Receive(socket.Conn, &msg)
			if err != nil {
				if err != io.EOF {
					socket.Logger.Debugf("Unexpected error: %s", err)
				}
				return
			}
			socket.Received <- msg
		}
	}()

	return socket
}

func (socket *WebSocket) Send(msg string) {
	err := websocket.Message.Send(socket.Conn, msg)
	if err != nil {
		socket.Logger.Error(err)
	}
}
