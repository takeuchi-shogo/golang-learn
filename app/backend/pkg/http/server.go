package http

import (
	"fmt"
	"net/http"
)

func NewServer(addr string, handler *http.ServeMux) error {
	if err := http.ListenAndServe(addr, handler); err != nil {
		return fmt.Errorf("failed to listen and serve: %w", err)
	}
	return nil
}
