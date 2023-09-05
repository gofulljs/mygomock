package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mux := InitMux()

	tests := []struct {
		name   string
		path   string
		before func(path string) *http.Request
	}{
		{
			name: "get hello",
			path: "/hello",
			before: func(path string) *http.Request {
				r, err := http.NewRequest("GET", path, nil)
				assert.NoError(t, err)

				return r
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			r := tt.before(tt.path)

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)

			assert.Equal(t, w.Code, http.StatusOK)
			assert.Equal(t, w.Body.String(), "Hello world!")
		})
	}

}
