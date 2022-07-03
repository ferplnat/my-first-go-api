package main

import (
	apifunctions "my-first-go-api/api-functions"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/albums", apifunctions.GetAlbums)
	router.POST("/albums", apifunctions.PostAlbums)

	router.Run("localhost:8080")
}
