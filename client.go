package main

import (
	"context"
	"time"
	"log"
	"github.com/danielholmes839/gRPC-chat/chat"
	"google.golang.org/grpc"
)

func main() {
	username := "Daniel"

	// Connect to the RPC service
	conn, _ := grpc.Dial(":4000", grpc.WithInsecure())
	defer conn.Close()

	//
	client := chat.NewChatServiceClient(conn)
	stream, err := client.Receive(context.Background(), &chat.Join{Username: username})
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages from the stream
	go func() {
		for {
			message, _ := stream.Recv()
			log.Printf("%s: %s", message.Username, message.Message)
		}
	}()

	// Send some messages
	go func() {
		for {
			client.Send(context.Background(), &chat.Message{Username: username, Message: "Hello"})
			time.Sleep(time.Second * 2)
		}
	}()
	time.Sleep(time.Minute)
}
