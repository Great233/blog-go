package test

import (
	"blog/config"
	"blog/models"
	"blog/pkg/response"
	"blog/router"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func init() {
	config.Init()
	models.Init()
	r = router.Init()
	response.Init()
}

var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7InVzZXJuYW1lIjoiYWRtaW4ifSwiYXVkIjoiR3JlYXQncyBibG9nIiwiZXhwIjoxNTg4NzQ4NTQ5LCJqdGkiOiJhNWZjYTgyMS02NTU5LTRmMjktYWNjYS0zM2M3ZjE3NzM5YzEiLCJpYXQiOjE1ODg3NDY3NDksImlzcyI6IkdyZWF0IiwibmJmIjoxNTg4NzQ2NzQ5LCJzdWIiOiJTZXNzaW9uIn0.JuG_ruXXRYapHxmDBQTkAFQp5Obo1ejGGfzvB18plEk"

func TestLogin(t *testing.T) {
	uri := "/web/login"
	resp := PostJson(uri, map[string]interface{}{
		"username": "admin",
		"password": "123456",
	}, nil, r)
	assert.Equal(t, 201, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body), resp.Header)
}

func TestWebGetArticles(t *testing.T) {
	uri := "/web/articles"
	resp := Get(uri, nil, r)
	assert.Equal(t, 401, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	resp = Get(uri, map[string]string{
		"xsrf-token": token,
	}, r)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func TestWebGetArticle(t *testing.T) {
	uri := "/web/articles/1"

	resp := Get(uri, map[string]string{
		"xsrf-token": token,
	}, r)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func TestWebGetTags(t *testing.T) {
	uri := "/web/tags"

	resp := Get(uri, map[string]string{
		"xsrf-token": token,
	}, r)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func TestWebGetTag(t *testing.T) {
	uri := "/web/tags/1"

	resp := Get(uri, map[string]string{
		"xsrf-token": token,
	}, r)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
