package tinyurl

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestCreateShortUrl(t *testing.T) {
	response := generateShortURL("hello")
	assert.Equal(t, len(response), 6)
}
