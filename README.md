# Simple gRPC chat

The simplest chat service I could make using gRPC.

```proto3
message Message {
    string username = 1;
    string message = 2;
}

message Join {
    string username = 1;
}

service ChatService {
    rpc Send(Message) returns (google.protobuf.Empty) {}
    rpc Receive(Join) returns (stream Message) {}
}```
