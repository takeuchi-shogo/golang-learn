package main

import (
	"fmt"
	"os"

	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/infra/routes"
	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/http"
)

func main() {
	fmt.Println("Hello, World!")

	err := http.NewServer("localhost:8080", routes.NewRouter())
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		os.Exit(1)
	}
}
