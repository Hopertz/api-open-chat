package main

import (
	"github/hopertz/api-open-chat/internal/websocket"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *application) ChatHandler(c *gin.Context) {
	w, r := c.Writer, c.Request

	websocket.ServeWs(app.pool, w, r)

}

func (app *application) ChatRoomHandler(c *gin.Context) {

	room_id := c.Param("id")

	var input struct {
		Page     int `form:"page"`
		PageSize int `form:"page_size" `
	}

	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if input.Page == 0 {
		input.Page = 1
	}

	if input.PageSize == 0 {
		input.PageSize = 20
	}

	id, err := strconv.Atoi(room_id)

	if err != nil {

		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	messages, err := app.models.MessageModel.FetchRoomMessages(id, input.Page, input.PageSize)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"messages": messages,
	})
}

func (app *application) ChatPrivateHandler(c *gin.Context) {

	user_id := c.Param("id")

	var input struct {
		Receiver int `form:"receiver"`
		Page     int `form:"page"`
		PageSize int `form:"page_size" `
	}

	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if input.Page == 0 {
		input.Page = 1
	}

	if input.PageSize == 0 {
		input.PageSize = 20
	}

	id, err := strconv.Atoi(user_id)

	if err != nil {

		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	messages, err := app.models.MessageModel.FetchPrivateMessages(id, input.Receiver, input.Page, input.PageSize)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"messages": messages,
	})
}
