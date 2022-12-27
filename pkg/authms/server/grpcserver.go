package server

import (
	"context"
	"net"

	"github.com/olegskip/messenger/pkg/authms/pb"
	"github.com/olegskip/messenger/pkg/authms/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	pb.UnimplementedAuthMsServer
	grpcServer *grpc.Server
	service service.IAuthService
}

func (g *GrpcServer) GetNewRt(ctx context.Context, req *pb.Credentials) (*pb.Tokens, error) {
	tokens, err := g.service.GetNewRt(
		service.CredentialsDto {
			Username: req.Username, 
			Password: req.Password,
		},
	)

	if err != nil {
		return nil, err
	}
	
	tokensPb := dtoToPbTokens(tokens)
	return &tokensPb, nil
}

func (g *GrpcServer) ExchangeRt(ctx context.Context, req *pb.Rt) (*pb.Tokens, error) {
	tokens, err := g.service.ExchangeRt(pbToDtoRt(req))

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tokensPb := dtoToPbTokens(tokens)
	return &tokensPb, nil
}

func (g *GrpcServer) RevokeRt(ctx context.Context, req *pb.Rt) (*emptypb.Empty, error) {
	return new(emptypb.Empty), g.service.RevokeRt(pbToDtoRt(req))
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

func dtoToPbRt(rtDto service.RtDto) (pb.Rt) {
	return pb.Rt {
		Token: rtDto.Token,
		ExpireTimestamp: timestamppb.New(rtDto.ExpireTimestamp),
	}
}

func dtoToPbJwt(jwtDto service.JwtDto) (pb.Jwt) {
	return pb.Jwt {
		Token: jwtDto.Token,
		ExpireTimestamp: timestamppb.New(jwtDto.ExpireTimestamp),
	}
}

func dtoToPbTokens(tokensDto service.TokensDto) (pb.Tokens) {
	rtPb := dtoToPbRt(tokensDto.Rt)
	jwtPb := dtoToPbJwt(tokensDto.Jwt)
	return pb.Tokens {
		Rt: &rtPb, 
		Jwt: &jwtPb,
	}
}

func pbToDtoRt(rtPb *pb.Rt) (rtDto service.RtDto) {
	return service.RtDto {
		Token: rtPb.Token,
		ExpireTimestamp: rtPb.ExpireTimestamp.AsTime(),
	}
}
