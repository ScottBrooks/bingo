package bingo

import "github.com/labstack/echo/v4"

// ServerInterface implemented by the handlers
type ServerInterface interface {
	Index(echo.Context) error
	Greeting(echo.Context) error
	BingoCard(echo.Context) error
	Pinger(echo.Context) error
	Memory(echo.Context) error
	Load(echo.Context) error
}

// RegisterHandlers register the handlers
func RegisterHandlers(router EchoRouter, si ServerInterface, m ...echo.MiddlewareFunc) {
	router.GET("/", si.Index)
	router.GET("/bingocard", si.BingoCard)
	router.GET("/greeting", si.Greeting)
	router.POST("/pinger", si.Pinger)
	router.GET("/memory", si.Memory)
	router.GET("/load", si.Load)
}
