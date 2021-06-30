package bingo

import "github.com/labstack/echo/v4"

// RegisterHandlers register the handlers
func RegisterHandlers(router EchoRouter, si *Hotwire, m ...echo.MiddlewareFunc) {
	router.GET("/", si.Index)                  // displays the main page
	router.GET("/card", si.Card)               // spawns a new bingo card
	router.GET("/cardsocket", si.CardSocket)   // spawns a new bingo card
	router.GET("/admin", si.Admin)             // admin view
	router.GET("/admincards", si.AdminCards)   // admin view
	router.GET("/adminevents", si.AdminEvents) // admin view
	router.POST("/adminevent", si.AdminEvent)  // admin view
	router.GET("/stats", si.Stats)             // admin view
}
