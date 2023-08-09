package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)


type Filters struct {
	Page     int
	PageSize int
}



func (app *application) readPageInt(c *gin.Context, key string) int {
	qs := c.DefaultQuery(key, "1")

	i, err := strconv.Atoi(qs)

	if err != nil {
		log.Println(err)
	}

	return i
}

func (app *application) readPageSizeInt(c *gin.Context, key string) int {
	qs := c.DefaultQuery(key, "20")

	i, err := strconv.Atoi(qs)

	if err != nil {
		log.Println(err)
	}

	return i
}
