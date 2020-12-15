package main

import (
	"context"
	"log"
	"time"

	"github.com/danielholmes839/gRPC-chat/chat"
	"google.golang.org/grpc"
)

func main() {
	username := "Daniel"

	// Connect to the RPC service and create a client
	conn, _ := grpc.Dial(":4000", grpc.WithInsecure())
	defer conn.Close()
	client := chat.NewChatServiceClient(conn)

	// Receive messages from the stream
	stream, _ := client.Receive(context.Background(), &chat.Join{Username: username})
	go func() {
		for {
			message, _ := stream.Recv()
			log.Printf("%s: %s", message.Username, message.Message)
		}
	}()
	time.Sleep(time.Second)

	// Send some messages
	go func() {
		for {
			client.Send(context.Background(), &chat.Message{Username: username, Message: "Hello"})
			time.Sleep(time.Second*5)
		}
	}()
	time.Sleep(time.Second * 30)
}
