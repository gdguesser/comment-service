// :gobuild e2e
package tests

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func createToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte("missionimpossible"))
	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}

func TestPostComment(t *testing.T) {
	t.Run("can post comment", func(t *testing.T) {
		client := resty.New()
		resp, err := client.R().
			SetHeader("Authorization", "bearer " + createToken()).
			SetBody(`{"slug": "/", "author": "Gabriel", "body": "hello peeps"}`).
			Post("http://localhost:8080/api/v1/comment")

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode())
	})

	t.Run("cannot post comment without JWT", func(t *testing.T) {
		client := resty.New()
		resp, err := client.R().
			SetBody(`{"slug": "/", "author": "Gabriel", "body": "hello peeps"}`).
			Post("http://localhost:8080/api/v1/comment")

		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode())
	})
}

func TestGetComment(t *testing.T) {
	t.Run("can get comment", func(t *testing.T) {
		client := resty.New()

		resp, err := client.R().Get("http://localhost:8080/api/v1/comment/01GEM0SMER7AK27A1NED7R17V6")
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode())
	})
}