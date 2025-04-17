package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "go-grpc-serviceA/generated"  // pb short for "protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
Service A:
- exposes: StatusService
- calls: GreeterService (on Service B)

Service B:
- exposes: GreeterService
- calls: StatusService (on Service A)
*/

// ========== StatusService Server (Service A) ==========
// define the statusServer struct to implement the gRPC server-side logic
type statusServer struct {
	pb.UnimplementedStatusServiceServer
	// UnimplementedStatusServiceServer is embedded to provide default method implementations (so your struct still satisfies the interface even if you don’t implement all RPCs yet).
}

// implement grpc handler (to handle grpc request from service B)
// this is server-side implementation of: rpc ReportStatus(StatusRequest) returns (StatusResponse);
func (s *statusServer) ReportStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	log.Printf("[Service A] Receive ReportStatus request from %s: %s", req.ServiceName, req.Status)
	return &pb.StatusResponse{Ack: "ACK from Service A"}, nil
}

// ========== Start Service A Server ==========
func startServiceAServer() {
	// binds gRPC server to TCP port 8081
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Service A: failed to start server on port :8081, %v", err)
	}

	// creates a new gRPC server instance
	// handles RPCs, manages lifecycles, supports middleware (interceptors)
	grpcServer := grpc.NewServer()

	// registers your implementation (&statusServer{}) with the gRPC server
	pb.RegisterStatusServiceServer(grpcServer, &statusServer{})
	log.Printf("Service A: Listening on :8081")

	// starts the server — blocks and listens for incoming gRPC calls.
	grpcServer.Serve(lis)
}

// ========== Call Service B ==========
// acting as a gRPC client calling method in service B
func callServiceB() {
	// establish a connection to service B
	// returns a grpc.ClientConn, which manages low-level HTTP/2 connections, connection pooling, and retries
	// insecure.NewCredentials() disables TLS (safe for local testing)
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Service A: failed to connect to service B, %v", err)
	}
	defer conn.Close()

	// create the gRPC stub
	// returns a client object that matches the GreeterService interface
	// This object knows how to encode requests, send them, wait, and decode responses
	client := pb.NewGreeterServiceClient(conn)

	// make the RPC call (call the SayHello method in service B)
	res, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "Greeting from Nguyen Anh Quan"})
	if err != nil {
		log.Fatalf("Service A: failed to call SayHello() in service B, %v", err)
	} else {
		log.Printf("Service A: Received response from service B: %v", res)
	}
}

func main() {
	// start service A server
	go startServiceAServer()

	time.Sleep(5 * time.Second)

	// call service B method
	log.Print("Service A: Start call method in service B")
	callServiceB()
}