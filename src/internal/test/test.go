package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
)

func MockHTTPServer(code int, payload string) (s *httptest.Server, err error) {
	body, err := os.ReadFile(payload)
	if err != nil {
		return s, fmt.Errorf("mock payload: %w", err)
	}
	s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Write(body)
	}))
	return
}
