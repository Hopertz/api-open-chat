package main

import (
	"github/hopertz/api-open-chat/internal/websocket"

	"github.com/gin-gonic/gin"
)

func (app *application) ChatHandler(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request

	websocket.ServeWs(app.pool, w, r)

}
