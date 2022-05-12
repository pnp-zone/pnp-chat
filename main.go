package main

import (
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/worker"
	"github.com/pnp-zone/common/conf"
	"golang.org/x/net/websocket"
	"gorm.io/gorm"
	"pnp-chat/broadcast"
)

var broadcastServer = broadcast.NewServer()

func chatHandler(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				c.Logger().Debug(err)
			}
		}()

		socket := NewWebSocket(c, ws)
		room := make(chan interface{}, 8)

		group := broadcastServer.GetGroup("test")
		group.Subscribe(room)
		defer group.Unsubscribe(room)

		for {
			select {
			case msg, ok := <-socket.Received:
				if ok {
					group.Publish(msg)
				} else {
					return
				}
			case msg := <-room:
				socket.Send(msg.(string))
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func MigrationHook() func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		return nil
	}
}

func RouterHook() func(e *echo.Echo, db *gorm.DB, config *conf.Config) error {
	return func(e *echo.Echo, db *gorm.DB, config *conf.Config) error {
		e.GET("chat/socket", chatHandler)
		return nil
	}
}

func StaticFileHook() (string, string) {
	return "main.js", "main.css"
}

func WorkerPoolHook() func(worker.Pool) error {
	return func(wp worker.Pool) error {
		go broadcastServer.Start()
		return nil
	}
}

var _ = MigrationHook
var _ = RouterHook
var _ = StaticFileHook
var _ = WorkerPoolHook
