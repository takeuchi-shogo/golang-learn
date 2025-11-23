package http

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func NewServer(addr string, handler *http.ServeMux) error {
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Printf("server is listening on %s", addr)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to listen and serve: %w", err)
	}
	return nil
}
