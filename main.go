package main

import (
	"log"
	"net"
	"reflect"

	hello_proto "github.com/WisdomEnigma/micro-peers/hello"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type helloServer struct {
	hello_proto.UnimplementedService
}

func NewServer() *helloServer { return &helloServer{} }

func (helloServer) Hello(resp *hello_proto.HelloRequest) (*hello_proto.HelloResponse, error) {

	if reflect.DeepEqual(resp, &hello_proto.HelloRequest{Message: ""}) {
		return &hello_proto.HelloResponse{Message: ""}, errors.Wrap(errors.New("user doesnot provide valid data"), "empty fields")
	}
	return &hello_proto.HelloResponse{Message: resp.Message}, nil
}

const Address = "127.0.0.1:9001"

func Server_init() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		return
	}

	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)

	hello_proto.RegisterService(server, NewServer())
	server.Serve(listen)

}
func main() {

	log.Println("Starting server...")
	Server_init()
}
