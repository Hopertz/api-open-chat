package main

import (
	"github/hopertz/api-open-chat/internal/websocket"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *application) ChatHandler(c *gin.Context) {
	w, r := c.Writer, c.Request

	websocket.ServeWs(app.pool, w, r)

}

func (app *application) ChatRoomHandler(c *gin.Context) {

	room_id := c.Param("id")

	id, err := strconv.Atoi(room_id)

	if err != nil {

		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	messages, err := app.models.MessageModel.FetchRoomMessages(id, 0)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"messages": messages,
	})
}
