package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	mux := InitMux()

	tests := []struct {
		name       string
		path       string
		before     func(path string) *http.Request
		wantResult any
	}{
		{
			name: "get hello",
			path: "/hello",
			before: func(path string) *http.Request {
				r, err := http.NewRequest("GET", path, nil)
				assert.NoError(t, err)

				return r
			},
			wantResult: "Hello, world!",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			r := tt.before(tt.path)

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, r)

			assert.Equal(t, w.Code, http.StatusOK)
			assert.Equal(t, tt.wantResult, w.Body.String())
		})
	}

}

// mock 请求的返回值
func TestMainMock(t *testing.T) {
	httpmock.Activate()

	tests := []struct {
		name          string
		path          string
		method        string
		initResponser func(status int, body []byte) (httpmock.Responder, error)
		wantStatus    int
		wantResult    []byte
	}{
		{
			name:   "mock string resp",
			path:   "/hello",
			method: "GET",
			initResponser: func(status int, body []byte) (httpmock.Responder, error) {
				return httpmock.NewStringResponder(status, string(body)), nil
			},
			wantStatus: http.StatusOK,
			wantResult: []byte("hello world"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "http://localhost" + tt.path

			responder, err := tt.initResponser(http.StatusOK, tt.wantResult)
			assert.NoError(t, err)

			httpmock.RegisterResponder(
				// 设置拦截的http方法
				tt.method,
				// 设置需要拦截的url
				url,
				// 设置需要替换成什么返回值
				responder,
			)

			req, err := http.NewRequest(tt.method, url, nil)
			assert.NoError(t, err)

			c := &http.Client{}
			res, err := c.Do(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			assert.Equal(t, res.StatusCode, http.StatusOK)
			assert.Equal(t, tt.wantResult, body)
		})
	}
}
