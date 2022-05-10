package main

import (
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/worker"
	"github.com/pnp-zone/common/conf"
	"golang.org/x/net/websocket"
	"gorm.io/gorm"
)

func chatHandler(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			if len(msg) > 0 {
				err = websocket.Message.Send(ws, msg)
				if err != nil {
					c.Logger().Error(err)
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
		return nil
	}
}
