package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	tdb "github.com/tinyurl/internal/db"
	"github.com/tinyurl/pkg/tinyurl"
)

func main() {
	db, err := tdb.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	defer db.Close()
	router := defineRouter(db)

	// Run the server
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

func defineRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	handler := tinyurl.Handler{
		Db: db,
	}
	// Define routes for fetching tiny urls
	router.GET("/", handler.GetUrls)
	router.GET("/:shortUrl", handler.GetUrl)
	router.GET("/:shortUrl/stats", handler.GetUrlStats)

	// Define routes for creating tiny urls
	router.POST("/", handler.PostUrl)

	// Define routes for deleting tiny urls
	router.DELETE("/:shortUrl", handler.DeleteUrl)
	return router
}
