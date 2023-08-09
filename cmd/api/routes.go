package main

import "github.com/gin-gonic/gin"

func (app *application) routes() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/ws", app.ChatHandler)

	r.GET("/rooms/:id", app.ChatRoomHandler)

	r.Use(enableCORS())

	return r
}
