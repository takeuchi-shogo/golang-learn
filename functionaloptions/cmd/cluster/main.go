package main

import (
	"context"
	"fmt"
	"log"

	"github.com/takeuchi-shogo/golang-learn/functionaloptions/cluster"
)

func main() {
	cluster, err := cluster.New(context.Background(), cluster.WithMaxWorkers(10))
	if err != nil {
		log.Fatalf("failed to create cluster: %v", err)
	}
	fmt.Println(cluster.GetMaxWorkers())
}
