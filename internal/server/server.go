package server

import (
	"fmt"
	"net/http"
)

func NewServer(h http.Handler, port int) *http.Server {

	s := &http.Server{
		Handler: h,
		Addr:    fmt.Sprintf("localhost:%d", port),
	}
	return s
}
