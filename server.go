package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func startServer() (*gin.Engine, error) {
	server := gin.Default()
	// api
	apiRoute := server.Group("/api")
	apiRoute.GET("/tunnels", tunnelHandle)
	// line callback
	lineRoute := server.Group("/line")
	{
		handler, err := BuildLinebotHandler()
		if err != nil {
			return nil, err
		}
		lineRoute.POST("/callback", WrapHttpHandlerToGin(handler))
	}

	return server, nil
}

func tunnelHandle(c *gin.Context) {
	if urls, err := ListTunnels(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    urls,
		})
	}
}

// convert http.HandlerFunc to gin's
func WrapHttpHandlerToGin(f http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(c.Writer, c.Request)
	}
}

func WrapHttpHandlerToGinBuilder(f http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(c.Writer, c.Request)
	}
}
