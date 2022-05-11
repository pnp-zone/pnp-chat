package main

import (
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/worker"
	"github.com/pnp-zone/common/conf"
	"golang.org/x/net/websocket"
	"gorm.io/gorm"
)

type ChatUser = WebSocket

func chatHandler(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				c.Logger().Debug(err)
			}
		}()

		user := NewWebSocket(c, ws)
		room := chatServer.GetChatRoom("test")
		room.Add(user)
		defer room.Remove(user)

		for {
			select {
			case msg, ok := <-user.Received:
				if ok {
					room.Send(msg)
				} else {
					return
				}
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
		go chatServer.HandleCalls()
		return nil
	}
}

var _ = MigrationHook
var _ = RouterHook
var _ = StaticFileHook
var _ = WorkerPoolHook
