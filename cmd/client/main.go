package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ZutrixPog/gterminal/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Cyan   = "\033[36m"
	Yellow = "\033[33m"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

// TODO: clean up the mess
func main() {
	flag.Parse()
	var opts []grpc.DialOption

	if *tls {
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	client := client.NewTerminalClient(*serverAddr, opts)
	defer client.Close()

	for {
		var username, password string
		fmt.Print(string(Yellow), "Enter username: ", string(Reset))
		fmt.Scanf("%s", &username)
		fmt.Print(string(Yellow), "Enter password: ", string(Reset))
		fmt.Scanf("%s", &password)

		if err := client.SignIn(username, password); err == nil {
			fmt.Println(string(Red), "Connected. Enter your commands:", string(Reset))
			break
		}
	}

	input := make(chan string)
	output := client.RunCommand(input)

	fmt.Print(string(Cyan), "GTerminal> ", string(Reset))
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "done" {
			return
		}

		input <- line

		for b := range output {
			if b == "EOC" {
				break
			}
			fmt.Print(b)
		}
		fmt.Fprint(os.Stdout, "\r \r")
		fmt.Println()
		fmt.Print(string(Cyan), "GTerminal> ", string(Reset))
	}
}
