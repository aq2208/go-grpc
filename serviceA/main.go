package main

/*
Service A:
- exposes: StatusService
- calls: GreeterService (on Service B)

Service B:
- exposes: GreeterService
- calls: StatusService (on Service A)
*/

import (
	"context"
	"log"
	"net"
	"time"

	pb "go-grpc-serviceA/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ========== StatusService Server (Service A) ==========
// implement statusServer to handle requests
type statusServer struct {
	pb.UnimplementedStatusServiceServer
}

// handle grpc request from service B
func (s *statusServer) ReportStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	log.Printf("[Service A] Receive ReportStatus request from %s: %s", req.ServiceName, req.Status)
	return &pb.StatusResponse{Ack: "ACK from Service A"}, nil
}

// ========== Call Service B ==========
func callServiceB() {
	// establish a tcp connection to service B
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Service A: failed to connect to service B, %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterServiceClient(conn)

	// call method in service B
	res, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "Greeting from Nguyen Anh Quan"})
	if err != nil {
		log.Fatalf("Service A: failed to call SayHello() in service B, %v", err)
	} else {
		log.Printf("Service A: Received response from service B: %v", res)
	}
}

// ========== Start Service A Server ==========
func startServiceAServer() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Service A: failed to start server on port :8081, %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStatusServiceServer(grpcServer, &statusServer{})
	log.Printf("Service A: Listening on :8081")
	grpcServer.Serve(lis)
}

func main() {
	// start service A server
	go startServiceAServer()

	time.Sleep(5 * time.Second)

	// call service B method
	log.Print("Service A: Start call method in service B")
	callServiceB()
}