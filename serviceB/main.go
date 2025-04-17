package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "go-grpc-serviceB/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ========== GreeterService Server ==========
type greeterServer struct {
	pb.UnimplementedGreeterServiceServer
}

func (s *greeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("[Service B] Hello received from: %s", req.Name)
	return &pb.HelloResponse{Message: "Hello back from B"}, nil
}

// ========== Call Service A ==========
func callServiceA() {
	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Service B: failed to connect to A: %v", err)
	}
	defer conn.Close()

	client := pb.NewStatusServiceClient(conn)

	res, err := client.ReportStatus(context.Background(), &pb.StatusRequest{
		ServiceName: "Service B",
		Status:      "Running",
	})
	if err != nil {
		log.Fatalf("Service B: failed to call A: %v", err)
	}
	log.Println("[Service B] Received from A:", res.Ack)
}

func main() {
	// Start Service B server
	go func() {
		lis, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatalf("Service B: failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterGreeterServiceServer(grpcServer, &greeterServer{})
		log.Println("[Service B] Listening on :8080")
		grpcServer.Serve(lis)
	}()

	// Call Service A
	time.Sleep(1 * time.Second)
	callServiceA()

	select {}
}