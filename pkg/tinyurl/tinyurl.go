package tinyurl

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Db *sql.DB
}

type TinyUrl struct {
	ShortUrl       string     `json:"short_url"`
	LongUrl        string     `json:"long_url"`
	ExpirationTime *time.Time `json:"expiration_date,omitempty"`
}

type Stats struct {
	HitsLast24Hours int `json:"hits_last_24_hours,omitempty"`
	HitsLastWeek    int `json:"hits_last_week,omitempty"`
	HitsAllTime     int `json:"hits_all_time,omitempty"`
}

type PostTinyUrl struct {
	LongUrl        string `json:"long_url" binding:"required"`
	ExpirationTime string `json:"expiration_date,omit_empty"`
	ShortUrl       string
}

// GetUrls fetches the existing resources for a given user.
// For now this returns nothing.
func (h *Handler) GetUrls(c *gin.Context) {
	urls, err := getDbResources(h.Db)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "error fetching data from DB")
	}
	c.IndentedJSON(http.StatusOK, urls)
}

// GetUrl fetches the tiny url resource based on the unique indentifier
// and redirects you to the long url. For now, it just returns some mock data.
func (h *Handler) GetUrl(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	resource, err := getDbResource(h.Db, shortUrl)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	err = addDbUrlHit(h.Db, shortUrl)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	redirectLink := "https://" + resource.LongUrl
	http.Redirect(c.Writer, c.Request, redirectLink, http.StatusFound)
}

// GetUrlStats fetches the tiny url resource based on the unique indentifier.
// For now, it just returns some mock data.
func (h *Handler) GetUrlStats(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	resource, err := getDbResourceStats(h.Db, shortUrl)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, resource)
}

// PostUrl creates a tiny url resource. The tiny url short version must be unique.
func (h *Handler) PostUrl(c *gin.Context) {
	var tinyUrl PostTinyUrl
	if err := c.BindJSON(&tinyUrl); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad Payload")
		return
	}
	x := TinyUrl{
		LongUrl: tinyUrl.LongUrl,
	}
	if tinyUrl.ExpirationTime != "" {
		timeExpired, err := time.Parse("2006-01-02", tinyUrl.ExpirationTime)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Bad timestamp")
			return
		}
		x.ExpirationTime = &timeExpired
	}
	x.ShortUrl = generateShortURL(tinyUrl.LongUrl)

	err := createDbResource(h.Db, x)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "DB error")
		return
	}
	c.IndentedJSON(http.StatusOK, x)
}

// DeleteUrl deletes the tiny url resource.
func (h *Handler) DeleteUrl(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	err := deleteDbResource(h.Db, shortUrl)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	// TODO Add DB handler to delete resource
	c.IndentedJSON(http.StatusOK, fmt.Sprintf("Resource with ID %s deleted", shortUrl))
}

const (
	lengOfUrl = 6
)

// generateShortUrl creates a short URL based on the long URL provided and
// the current time string. One limitation is that it may return all
// digits.
func generateShortURL(longURL string) string {
	// Append the current timestamp to the long URL to make a random string.
	input := longURL + time.Now().String()
	hash := sha256.New()

	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	// Take a substring of the hash string to get the random string
	randomString := hashString[:lengOfUrl]

	return randomString
}
