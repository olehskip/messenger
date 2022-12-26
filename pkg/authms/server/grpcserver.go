package server

import (
	"context"
	"net"

	"github.com/olegskip/messenger/pkg/authms/pb"
	"github.com/olegskip/messenger/pkg/authms/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb.UnimplementedAuthMsServer
	grpcServer *grpc.Server
	service service.IAuthService
}

func (g *GrpcServer) GetNewRt(ctx context.Context, req *pb.Credentials) (*pb.Tokens, error) {
	rt, jwt, err := g.service.GetNewRt(
		service.CredentialsDto {
			Username: req.Username, 
			Password: req.Password,
		},
	)

	if err != nil {
		return nil, err
	}

	return &pb.Tokens{Jwt: &pb.Jwt{Token: jwt.Token}, Rt: &pb.Rt{Token: rt.Token}}, nil
}

func (g *GrpcServer) ExchangeRt(ctx context.Context, req *pb.Rt) (*pb.Rt, error) {
	rt, err := g.service.ExchangeRt(
		service.RtDto {
			Token: req.Token, 
		},
	)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.Rt{Token: rt.Token}, nil
}

func (g *GrpcServer) Run() error {
	lis, _ := net.Listen("tcp", "localhost:9092")
	if err := g.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func NewGRPCServer(service service.IAuthService) *GrpcServer {
	ans := GrpcServer{service: service}
	ans.grpcServer = grpc.NewServer()
	pb.RegisterAuthMsServer(ans.grpcServer, &ans)
	
	return &ans
}

