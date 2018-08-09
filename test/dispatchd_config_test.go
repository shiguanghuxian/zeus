package main

import (
	"encoding/json"
	"testing"

	"53it.net/zues/dispatchd"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:3200"
)

func TestDispatchdSayDispatchConfig(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("grpc server没有开启服务: %v", err)
	}
	defer conn.Close()
	c := dispatchd.NewDispatchConfigServiceClient(conn)

	r, err := c.SayDispatchConfig(context.Background(), &dispatchd.DispatchRequest{Topics: "test"})

	if err != nil {
		t.Fatalf("发消息失败: %v", err)
	}
	b, _ := json.Marshal(r)
	t.Log(string(b))
}
