package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func Get(uri string, headers map[string]string, router *gin.Engine) *http.Response {
	request := httptest.NewRequest("GET", uri, nil)
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)
	response := w.Result()
	defer response.Body.Close()
	return response
}

func PostJson(uri string, body map[string]interface{}, headers map[string]string, router *gin.Engine) *http.Response {
	params, _ := json.Marshal(body)
	request := httptest.NewRequest("POST", uri, bytes.NewReader(params))
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)
	response := w.Result()
	defer response.Body.Close()
	return response
}
