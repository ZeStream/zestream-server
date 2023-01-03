package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

// router is a function that returns a new instance of the Gin router with a single route defined:
func router() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	publicRoutes := router.Group("/api/v1")

	// a GET request to "/api/v1/ping" will execute the Ping function.
	publicRoutes.GET("/ping", Ping)

	return router
}

/*
MakeRequest is a helper function that creates an HTTP request and sends it to the router,
returning the response that the router returns.
It takes in three arguments:
  - method: the HTTP method to use (e.g. "GET", "POST")
  - url: the URL to send the request to
  - body: the request body, which will be JSON-encoded
*/
func makeRequest(method, url string, body interface{}) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	writer := httptest.NewRecorder()
	router().ServeHTTP(writer, request)
	return writer
}
