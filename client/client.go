package client

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/ZutrixPog/gterminal/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TerminalClient struct {
	conn   *grpc.ClientConn
	client *pb.TerminalClient
	ctx    context.Context
	token  string
}

func NewTerminalClient(address string, opts []grpc.DialOption) *TerminalClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	client := pb.NewTerminalClient(conn)

	ctx := context.Background()

	return &TerminalClient{
		conn:   conn,
		client: &client,
		ctx:    ctx,
		token:  "",
	}
}

func (t *TerminalClient) SignIn(username, password string) error {
	token, err := (*t.client).SignIn(t.ctx, &pb.Credentials{Username: username, Password: password})
	if err != nil {
		return err
	}

	t.token = token.Token
	return nil
}

func (t *TerminalClient) RunCommand(commands chan string) <-chan string {
	out := make(chan string)
	if t.token == "" {
		return nil
	}

	stream, err := (*t.client).Run(t.ctx)
	if err != nil {
		fmt.Println("couldnt connect")
		return nil
	}

	go t.handleSend(&stream, commands)

	go t.handleReceive(&stream, out)

	go t.handleExit(out)

	return out
}

func (t *TerminalClient) handleSend(stream *pb.Terminal_RunClient, commands chan string) {
	for command := range commands {
		if err := (*stream).Send(&pb.Command{Text: command, Token: t.token}); err != nil {
			log.Fatalf("can not send %v", err)
		}
	}
}

func (t *TerminalClient) handleReceive(stream *pb.Terminal_RunClient, out chan string) {
	for {
		res, err := (*stream).Recv()
		if err == io.EOF {
			close(out)
			return
		}
		if err != nil {
			//log.Fatalf("can not receive %v", err)
			log.Fatal("Disconnected")
		}
		out <- res.Text
	}
}

func (t *TerminalClient) handleExit(out chan string) {
	<-t.ctx.Done()
	if err := t.ctx.Err(); err != nil {
		log.Println(err)
	}
	close(out)
}

func (t *TerminalClient) Close() {
	t.conn.Close()
}
