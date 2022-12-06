package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/ZutrixPog/gterminal/grpc"
	"github.com/ZutrixPog/gterminal/server"
	"github.com/ZutrixPog/gterminal/terminal"
	"google.golang.org/grpc"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	userDbFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	terminals := terminal.NewMemTerminalRepo()
	server := server.NewCommandServer(terminals)

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTerminalServer(grpcServer, server)
	fmt.Println("Server is running")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to start server")
	}
}
