package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	e := echo.New()
	e.Use(loginMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, helloWorldHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		res := rec.Body.String()

		assert.Equal(t, res, "Hello, World!")
	}
}
