package main

import (
	"github.com/olegskip/messenger/pkg/authms/server"
	"github.com/olegskip/messenger/pkg/authms/service"
	"github.com/olegskip/messenger/pkg/authms/dal"

	"log"
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

type Vehicle struct {
}

type Car struct {
    Vehicle //anonymous field Vehicle
}

func test(v *Vehicle) {

}

func main() {
	if err := server.NewGRPCServer(
		service.NewAuthService(
			new(dal.InMemoryDao),
		),
	).Run(); err != nil {
		log.Fatalf("Can't run server; Error = %v", err)
	}
}
