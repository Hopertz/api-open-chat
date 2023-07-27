package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func enableCORS() gin.HandlerFunc {

	corsConfig := cors.DefaultConfig()

	corsConfig.AllowAllOrigins = true

	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}

	return cors.New(corsConfig)

}
