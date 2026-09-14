package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	router := gin.New()

	prepareHandlers(router)

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/ping", nil)
	require.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	m.Run()
}
