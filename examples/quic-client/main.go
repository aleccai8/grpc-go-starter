package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"github.com/zhengheng7913/grpc-go-starter/pkg/codec"
	_ "github.com/zhengheng7913/grpc-go-starter/plugin/grpc-protocol-plugin"
	"github.com/zhengheng7913/grpc-go-starter/plugin/grpc-protocol-plugin/stream"
)

type Message struct {
	Payload string
}

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic"},
	}
	conn, err := quic.DialAddr("127.0.0.1:8090", tlsConf, nil)
	if err != nil {
		panic(err)
	}

	s, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		panic(err)
	}
	cs := stream.NewStream(s, s, codec.GetClient("json"))
	send := &Message{
		Payload: "hello world",
	}
	err = cs.SendMsg(send)
	if err != nil {
		panic(err)
	}

	for {
		recv := &Message{}
		err = cs.RecvMsg(recv)
		if err != nil {
			panic(err)
		}
		fmt.Println("recv:", recv.Payload)
	}
	s.Close()
}
