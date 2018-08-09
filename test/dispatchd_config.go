package main

import (
	"fmt"

	"53it.net/zues/dispatchd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:3200"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("grpc server没有开启服务: %v", err)
	}
	defer conn.Close()
	c := dispatchd.NewDispatchConfigServiceClient(conn)

	r, err := c.SayDispatchConfig(context.Background(), &dispatchd.DispatchRequest{})
	fmt.Println(r, err)
}
