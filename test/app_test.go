package test

import (
	"blog/config"
	"blog/models"
	"blog/router"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)
var r *gin.Engine
func init() {
	config.Init()
	models.Init()
	r = router.Init()
}

func TestGetArticles(t *testing.T) {
	uri := "/app/articles"
	response := Get(uri, nil, r)
	assert.Equal(t, 200, response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	log.Println(string(body))

	uri = "/app/articles/a"
	response = Get(uri, nil, r)
	assert.Equal(t, 404, response.StatusCode)
	body, _ = ioutil.ReadAll(response.Body)
	log.Println(string(body))

	uri = "/app/articles/aa"
	response = Get(uri, nil, r)
	assert.Equal(t, 200, response.StatusCode)
	body, _ = ioutil.ReadAll(response.Body)
	log.Println(string(body))
}

func TestGetArticle(t *testing.T) {
	uri := "/app/articles/aa"
	response := Get(uri, nil, r)
	assert.Equal(t, 200, response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	log.Println(string(body))

	uri = "/app/articles/bb"
	response = Get(uri, nil, r)
	assert.Equal(t, 404, response.StatusCode)
	body, _ = ioutil.ReadAll(response.Body)
	log.Println(string(body))

	uri = "/app/articles/bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	response = Get(uri, nil, r)
	assert.Equal(t, 404, response.StatusCode)
	body, _ = ioutil.ReadAll(response.Body)
	log.Println(string(body))
}

func TestGetTags(t *testing.T) {
	uri := "/app/tags"
	response := Get(uri, nil, r)
	assert.Equal(t, 200, response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	log.Println(string(body))
}
