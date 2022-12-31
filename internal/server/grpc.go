package server

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/multimoml/dispatcher/internal/proto"
)

func GetConfig(key string) (string, error) {
	conn, err := grpc.Dial(sConfig.ConfigServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	log.Println("Connecting to config server at", sConfig.ConfigServer)

	if err != nil {
		return "", err
	}
	defer func(conn *grpc.ClientConn) {
		if conn.Close() != nil {
			log.Println("Error closing connection to config server:", err)
		}
	}(conn)

	client := proto.NewConfigClient(conn)
	value, err := client.GetConfig(context.Background(), &proto.ConfigRequest{Key: key})

	if err != nil {
		return "", err
	}

	return value.Value, nil
}
