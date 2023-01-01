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
	// "google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	pb.UnimplementedAuthMsServer
	grpcServer *grpc.Server
	service service.IAuthService
}

func (g *GrpcServer) GetNewRefreshToken(ctx context.Context, req *pb.Credentials) (*pb.Tokens, error) {
	tokens, err := g.service.GetNewRefreshToken(
		service.CredentialsDto {
			UserUuid: req.Uuid, 
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
	tokens, err := g.service.ExchangeRefreshToken(req.HashedRefreshToken)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tokensPb := dtoToPbTokens(tokens)
	return &tokensPb, nil
}

func (g *GrpcServer) RevokeRefreshToken(ctx context.Context, req *pb.RefreshToken) (*emptypb.Empty, error) {
	return new(emptypb.Empty), g.service.RevokeRefreshToken(req.HashedRefreshToken)
}

func (g *GrpcServer) GetTokenOwner(ctx context.Context, req *pb.AccessToken) (*pb.TokenOwner, error) {
	uuid, isTokenRevoked, err := g.service.GetUserUuid(req.HashedAccessToken)

	if err != nil {
		return nil, err
	}

	return &pb.TokenOwner{Uuid: uuid, IsTokenRevoked: isTokenRevoked}, nil
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

// func dtoToPbRefreshToken(refreshTokenDto service.RefreshTokenDto) (pb.RefreshToken) {
// 	return pb.RefreshToken {
// 		Token: refreshTokenDto.Token,
// 		ExpireTimestamp: timestamppb.New(refreshTokenDto.ExpireTimestamp),
// 	}
// }
//
// func dtoToPbAccessToken(accessTokenDto service.AccessTokenDto) (pb.AccessToken) {
// 	return pb.AccessToken {
// 		Token: accessTokenDto.Token,
// 		ExpireTimestamp: timestamppb.New(accessTokenDto.ExpireTimestamp),
// 	}
// }

func dtoToPbTokens(hTokensDto service.HTokensPairDto) (pb.Tokens) {
	// refreshTokenPb := dtoToPbRefreshToken(tokensDto.RefreshToken)
	// accessTokenPb := dtoToPbAccessToken(tokensDto.AccessToken)
	return pb.Tokens {
		HashedRefreshToken: hTokensDto.HashedRefreshToken, 
		HashedAccessToken: hTokensDto.HashedAccessToken,
	}
}
//
// func pbToDtoRefreshToken(refreshTokenPb *pb.RefreshToken) (rtDto service.RefreshTokenDto) {
// 	return service.RefreshTokenDto {
// 		Token: refreshTokenPb.Token,
// 		ExpireTimestamp: refreshTokenPb.ExpireTimestamp.AsTime(),
// 	}
// }
