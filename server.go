package main

import (
	"fmt"
	"log"
	"net"

	"github.com/danielholmes839/gRPC-chat/chat"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// RPCs type
type RPCs map[string]chat.ChatService_ReceiveServer

// ChatServer struct
type ChatServer struct {
	clients RPCs
}

// Broadcast func
func (s *ChatServer) Broadcast(message *chat.Message) {
	for username := range s.clients {
		s.clients[username].Send(message)
	}
}

// Send function (RPC)
func (s *ChatServer) Send(ctx context.Context, message *chat.Message) (*emptypb.Empty, error) {
	log.Printf("New message from %s: %s", message.Username, message.Message)
	s.Broadcast(message)
	return &emptypb.Empty{}, nil
}

// Receive function (RPC)
func (s *ChatServer) Receive(join *chat.Join, client chat.ChatService_ReceiveServer) error {
	// RPC connects
	s.clients[join.Username] = client
	connected := fmt.Sprintf("%s has joined the server!", join.Username)
	log.Println(connected)
	s.Broadcast(&chat.Message{Username: "SERVER", Message: connected})
	<-client.Context().Done()

	// RPC disconnects
	delete(s.clients, join.Username)
	disconnected := fmt.Sprintf("%s has disconnected", join.Username)
	log.Println(disconnected)
	s.Broadcast(&chat.Message{Username: "SERVER", Message: disconnected})
	return nil
}

func main() {
	fmt.Println("Listening...")
	lis, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	chat.RegisterChatServiceServer(server, &ChatServer{make(RPCs)})
	server.Serve(lis)
}
