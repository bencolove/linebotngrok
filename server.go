package main

import (
	"net/http"

	"com.roger.ngrok.linebot/ngrokapi"

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
		lineRoute.POST("/callback", WrapHttpHandlerFuncToGinBuilder(handler))
	}
	// graphql
	graphqlRoute := server.Group("/g")
	{
		graphqlRoute.GET("/users", WrapHttpHandlerFuncToGinBuilder(ngrokapi.GetNgrokUsers))
		graphqlRoute.POST("/", WrapHttpHandlerToGin(ngrokapi.GetGraphqlHttpHandler()))
		graphqlRoute.GET("/tunnels", WrapHttpHandlerFuncToGinBuilder(ngrokapi.GetNgrokTunnels))
	}

	return server, nil
}

func tunnelHandle(c *gin.Context) {
	if urls, err := ngrokapi.ListTunnels(); err != nil {
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
func WrapHttpHandlerToGin(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func WrapHttpHandlerFuncToGinBuilder(f http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(c.Writer, c.Request)
	}
}
