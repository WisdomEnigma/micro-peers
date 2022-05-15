package hello_proto

import (
	"context"
	"log"

	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// grpc supported versions
const _ = grpc.SupportPackageIsVersion7

type Hello_Service interface {

	// 	Hello Service is a service that will discover when connect with network mesh
	Hello(resp *HelloRequest, callback ...grpc.CallOption) (*HelloResponse, error)
}

type HelloClient struct {

	// grpc client connection state
	Client grpc.ClientConnInterface

	// grpc timeout connection
	Ctx context.Context
}

func NewClient(ctx context.Context, client grpc.ClientConnInterface) Hello_Service {
	return &HelloClient{Ctx: ctx, Client: client}
}

func (c *HelloClient) Hello(resp *HelloRequest, callback ...grpc.CallOption) (*HelloResponse, error) {

	out := new(HelloRequest)

	config := consul.DefaultConfig()

	_consulClient, err := consul.NewClient(config)
	if err != nil {
		log.Fatalln("Consul connection error: ", err)
		return &HelloResponse{Message: ""}, err
	}

	kv := _consulClient.KV()
	pairs := consul.KVPair{Key: "grpc_cli", Value: []byte(resp.Message)}

	err = c.Client.Invoke(c.Ctx, "/hello.HelloService/Hello", resp, out, callback...)
	if err != nil {
		return &HelloResponse{Message: ""}, err
	}

	_, err = kv.Put(&pairs, nil)
	if err != nil {
		log.Fatalln("Consul reject request", err)
		return &HelloResponse{Message: ""}, err
	}

	return &HelloResponse{Message: resp.Message}, nil

}

type Hello_Server interface {
	Hello(resp *HelloRequest) (*HelloResponse, error)
	mustEmbedUnimplemented()
}

// type Hello_Service_Server struct{}

type UnimplementedService struct{}

func (UnimplementedService) Hello(resp *HelloRequest) (*HelloResponse, error) {
	return &HelloResponse{Message: ""}, status.Errorf(codes.Unimplemented, "hello request is implement")
}

func (UnimplementedService) mustEmbedUnimplemented() {}

type UnSafeTaskService interface {
	mustEmbedUnimplemented()
}

func RegisterService(register grpc.ServiceRegistrar, sv Hello_Server) {
	register.RegisterService(&Hello_Desc, sv)
}

func _Hello_Service_Handler(s interface{}, ctx context.Context, desc func(interface{}) error, intercept grpc.UnaryServerInterceptor) (interface{}, error) {

	in := new(HelloRequest)

	err := desc(in)
	if err != nil {
		return nil, err
	}

	if intercept == nil {
		return s.(Hello_Server).Hello(in)
	}

	unary := &grpc.UnaryServerInfo{

		Server:     s,
		FullMethod: "hello.HelloService/Hello",
	}

	handler := func(ctx context.Context, resp interface{}) (interface{}, error) {

		req, err := s.(Hello_Server).Hello(in)
		return req, err
	}

	return intercept(ctx, in, unary, handler)
}

var Hello_Desc = grpc.ServiceDesc{
	ServiceName: "hello.HelloService",
	HandlerType: (*Hello_Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _Hello_Service_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/WisdomEnigma/micro-peers/hello/hello.proto",
}

const port = ":9001"

func HelloCLientInit() {

	log.Println("Starting .....")

	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

		log.Fatalln("Connection Failed: ", err)
		return
	}

	defer conn.Close()

	req, err := NewClient(context.Background(), conn).Hello(&HelloRequest{Message: "hello"})

	if err != nil {

		log.Fatalln("Client throw :", err)
		return
	}

	log.Println("Response:", req.Message)

}
