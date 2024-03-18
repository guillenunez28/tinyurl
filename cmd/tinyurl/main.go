package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinyurl/pkg/tinyurl"
)

func main() {
	router := gin.Default()

	// Define routes for fetching tiny urls
	router.GET("/", tinyurl.GetUrls)
	router.GET("/:shortUrl", tinyurl.GetUrl)
	router.GET("/:shortUrl/stats", tinyurl.GetUrlStats)

	// Define routes for creating tiny urls
	router.POST("/", tinyurl.PostUrl)

	// Define routes for deleting tiny urls
	router.DELETE("/:id", tinyurl.DeleteUrl)

	// Run the server
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
