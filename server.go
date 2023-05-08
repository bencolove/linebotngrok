package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func startServer() *gin.Engine {
	server := gin.Default()
	server.GET("/tunnels", tunnelHandle)
	return server
}

func tunnelHandle(c *gin.Context) {
	if urls, err := getTunnels(apiKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    urls,
		})
	}
}

func getTunnels(apiKey string) ([]string, error) {
	tunnelResp, err := fetchTunnels(apiKey)
	if err != nil {
		return nil, err
	} else {
		fmt.Printf("tunnels: %+v\n", tunnelResp)
		urls := []string{}
		if len(tunnelResp.Tunnels) > 0 {
			url := tunnelResp.Tunnels[0].PublicURL
			if url != "" {
				urls = append(urls, url)
			}
		}
		return urls, nil
	}

}
