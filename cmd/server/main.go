package main

import (
	"context"

	"github.com/multimoml/dispatcher/internal/server"
)

func main() {
	ctx := context.Background()
	server.Run(ctx)
}
