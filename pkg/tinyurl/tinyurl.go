package tinyurl

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TinyUrl struct {
	ShortVersion    string    `json:"short_version"`
	LongVersion     string    `json:"long_version"`
	ExpirationTime  time.Time `json:"expiration_time"`
	HitsLast24Hours int       `json:"hits_last_24_hours"`
	HitsLastWeek    int       `json:"hits_last_week"`
	HitsAllTime     int       `json:"hits_all_time"`
}

type PostTinyUrl struct {
	LongVersion    string    `json:"long_version" binding:"required"`
	ExpirationTime time.Time `json:"expiration_time"`
}

// GetUrls fetches the existing resources for a given user.
// For now this returns nothing.
func GetUrls(c *gin.Context) {
	// TODO Add DB handler to fetch the resources
	c.IndentedJSON(http.StatusOK, []string{})
}

// GetUrl fetches the tiny url resource based on the unique indentifier
// and redirects you to the long url. For now, it just returns some mock data.
func GetUrl(c *gin.Context) {
	// TODO Add DB handler to fetch the resource
	shortUrl := c.Query("shortUrl")
	// TODO Add DB handler to fetch the resource
	resource := TinyUrl{
		ShortVersion:   shortUrl,
		LongVersion:    "https://www.google.com",
		ExpirationTime: time.Now(),
	}
	http.Redirect(c.Writer, c.Request, resource.LongVersion, http.StatusFound)
}

// GetUrlStats fetches the tiny url resource based on the unique indentifier.
// For now, it just returns some mock data.
func GetUrlStats(c *gin.Context) {
	// TODO Add DB handler to fetch the resource
	shortUrl := c.Query("shortUrl")
	// TODO Add DB handler to fetch the resource
	resource := TinyUrl{
		ShortVersion:    shortUrl,
		LongVersion:     "www.google.com",
		ExpirationTime:  time.Now(),
		HitsAllTime:     4,
		HitsLast24Hours: 1,
		HitsLastWeek:    2,
	}
	c.IndentedJSON(http.StatusOK, resource)
}

// PostUrl creates a tiny url resource. The tiny url short version must be unique.
func PostUrl(c *gin.Context) {
	var tinyUrl PostTinyUrl
	if err := c.BindJSON(&tinyUrl); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad Payload")
		return
	}
	c.IndentedJSON(http.StatusOK, nil)
}

// DeleteUrl deletes the tiny url resource.
func DeleteUrl(c *gin.Context) {
	id := c.Query("id")
	// TODO Add DB handler to delete resource
	c.IndentedJSON(http.StatusOK, fmt.Sprintf("Resource with ID %s deleted", id))
}
