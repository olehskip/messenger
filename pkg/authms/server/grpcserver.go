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

func (g *GrpcServer) GetNewRefreshToken(ctx context.Context, req *pb.Credentials) (*pb.Tokens, error) {
	tokens, err := g.service.GetNewRefreshToken(
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

func (g *GrpcServer) ExchangeRefreshToken(ctx context.Context, req *pb.RefreshToken) (*pb.Tokens, error) {
	tokens, err := g.service.ExchangeRefreshToken(pbToDtoRefreshToken(req))

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tokensPb := dtoToPbTokens(tokens)
	return &tokensPb, nil
}

func (g *GrpcServer) RevokeRefreshToken(ctx context.Context, req *pb.RefreshToken) (*emptypb.Empty, error) {
	return new(emptypb.Empty), g.service.RevokeRefreshToken(pbToDtoRefreshToken(req))
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

func dtoToPbRefreshToken(refreshTokenDto service.RefreshTokenDto) (pb.RefreshToken) {
	return pb.RefreshToken {
		Token: refreshTokenDto.Token,
		ExpireTimestamp: timestamppb.New(refreshTokenDto.ExpireTimestamp),
	}
}

func dtoToPbAccessToken(accessTokenDto service.AccessTokenDto) (pb.AccessToken) {
	return pb.AccessToken {
		Token: accessTokenDto.Token,
		ExpireTimestamp: timestamppb.New(accessTokenDto.ExpireTimestamp),
	}
}

func dtoToPbTokens(tokensDto service.TokensDto) (pb.Tokens) {
	refreshTokenPb := dtoToPbRefreshToken(tokensDto.RefreshToken)
	accessTokenPb := dtoToPbAccessToken(tokensDto.AccessToken)
	return pb.Tokens {
		RefreshToken: &refreshTokenPb, 
		AccessToken: &accessTokenPb,
	}
}

func pbToDtoRefreshToken(refreshTokenPb *pb.RefreshToken) (rtDto service.RefreshTokenDto) {
	return service.RefreshTokenDto {
		Token: refreshTokenPb.Token,
		ExpireTimestamp: refreshTokenPb.ExpireTimestamp.AsTime(),
	}
}
