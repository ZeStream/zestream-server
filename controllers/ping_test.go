package controllers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

/*
TestPing is a test function that sends a GET request to the "/api/v1/ping" route and
verifies that the response has a 200 OK status code and a body of "pong".
*/
func TestPing(t *testing.T) {
	recorder := makeRequest("GET", "/api/v1/ping", nil)

	// Assert that the response status is 200 OK
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Assert that the response body is "pong"
	assert.Equal(t, "pong", recorder.Body.String())
}
