package transport

import (
	"context"
	"fmt"

	grpcsvr "github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	
	"github.com/olegskip/messenger/pkg/authms/pb"
)

var service micro.Service
type ITransport interface {
	Auth(ctx context.Context, req *pb.Credentials, rsp *pb.LoginResponse) error

}

type AuthCredentials struct {
	Username string
	Password string
}

type AuthResponse struct {
	Token string
}

type GRPCTransport struct {
}


func CreateGRPCTransport() *GRPCTransport {
	service = micro.NewService(
		micro.Server(grpcsvr.NewServer()),
		micro.Name("AuthMs"),
		micro.Address("localhost:50052"),
	)

	service.Init()
	transportRes := GRPCTransport{}
	pb.RegisterAuthMsHandler(service.Server(), &transportRes)

	return &transportRes
}

func Run() {	
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func (g *GRPCTransport) Auth(ctx context.Context, req *pb.Credentials, rsp *pb.LoginResponse) error {
	// _, err := authMs.Auth(req.Username, req.Password)
	// rsp.IsSuccess = err == nil

	return nil
}

