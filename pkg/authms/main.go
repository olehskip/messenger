package main

import (
	"github.com/olegskip/messenger/pkg/authms/transport"
)


// Setup and the client
// func runClient(service micro.Service) {
// 	greeter := pb.NewGreeterService("Greeter", service.Client())
//
// 	// Call the greeter
// 	rsp, err := greeter.Hello(context.TODO(), &pb.Request{Name: "John1"})
// 	if err != nil {
// 		fmt.Println("ERROR!")
// 		fmt.Println(err)
// 		return
// 	}
//
// 	fmt.Println(rsp.Greeting)
// }

func main() {
	// var iTransport transport.ITransport

	transport.CreateGRPCTransport()
	transport.Run()
}
