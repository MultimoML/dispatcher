package main

import (
	"context"

	"github.com/multimoml/dispatcher/internal/server"
)

func main() {
	server.Run(context.Background())
}
