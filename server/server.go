package server

import (
	"context"
	"fmt"
	"io"

	pb "github.com/ZutrixPog/gterminal/grpc"
	"github.com/ZutrixPog/gterminal/terminal"
	"github.com/ZutrixPog/gterminal/utils"
)

type TerminalServer struct {
	pb.UnimplementedTerminalServer

	terminals terminal.TerminalRepo
}

func NewCommandServer(terminals terminal.TerminalRepo) *TerminalServer {
	return &TerminalServer{
		terminals: terminals,
	}
}

func (c *TerminalServer) Run(stream pb.Terminal_RunServer) error {
	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Println(err)
			return err
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := utils.VerifyToken(in.Token); err != nil {
			return err
		}

		terminal := c.terminals.GetTerminal(in.Token)

		out, err := terminal.Run(in.Text)
		if err != nil {
			stream.Send(&pb.Command{Text: "Command not found", Token: in.Token})
			stream.Send(&pb.Command{Text: "EOC", Token: in.Token})
			continue
		}

		if out != nil {
			b := make([]byte, 8)
			for {
				_, err := out.Read(b)
				if err == io.EOF {
					break
				}
				stream.Send(&pb.Command{Text: string(b), Token: in.Token})
			}
		}
		stream.Send(&pb.Command{Text: "EOC", Token: in.Token})
	}
}

// Todo: Add real authentication
func (c *TerminalServer) SignIn(ctx context.Context, creds *pb.Credentials) (*pb.Token, error) {
	token, err := utils.GenerateToken(creds.Username)
	if err != nil {
		return nil, err
	}

	c.terminals.SetTerminal(token, terminal.NewTerminal("/"))

	return &pb.Token{Token: token}, nil
}
