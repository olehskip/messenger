package server

import (
	"context"
	"net"

	// "fmt"
	"github.com/olegskip/messenger/pkg/authms/pb"
	"github.com/olegskip/messenger/pkg/authms/service"
	"google.golang.org/grpc"
)

type AuthCredentials struct {
	Username string
	Password string
}

type AuthResponse struct {
	Token string
}

type IServer interface {
	GetNewRT(ctx context.Context, req *pb.Credentials) (*pb.Tokens, error)
	Run() error
}

type GRPCServer struct {
	pb.UnimplementedAuthMsServer
	grpcServer *grpc.Server
	service service.IAuthService
}

func (g *GRPCServer) GetNewRT(ctx context.Context, req *pb.Credentials) (*pb.Tokens, error) {
	rt, jwt, _ := g.service.GetNewRT(req.Username, req.Password)
	
	return &pb.Tokens{Jwt: jwt, Rt: rt}, nil
}

func (g *GRPCServer) Run() error {
	lis, _ := net.Listen("tcp", "localhost:9092")
	if err := g.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
func NewGRPCServer(service service.IAuthService) *GRPCServer {
	ans := GRPCServer{service: service}
	ans.grpcServer = grpc.NewServer()
	pb.RegisterAuthMsServer(ans.grpcServer, &ans)
	
	return &ans
}

