package main

import (
	"context"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"github.com/zhengheng7913/grpc-go-starter/api"
	"github.com/zhengheng7913/grpc-go-starter/pkg/codec"
	proto "github.com/zhengheng7913/grpc-go-starter/plugin/grpc-protocol-plugin"
	"github.com/zhengheng7913/grpc-go-starter/plugin/grpc-protocol-plugin/stream"
)

func newQuicServer() proto.Server {
	return &QuicServer{}
}

type Message struct {
	Payload string
}

type QuicServer struct {
}

func (q QuicServer) Serve(conn quic.Connection) error {
	s, err := conn.AcceptStream(context.Background())
	ss := stream.NewStream(s, s, codec.GetServer("json"))
	if err != nil {
		panic(err)
	}
	for {
		recv := &Message{}
		err = ss.RecvMsg(recv)
		if err != nil {
			panic(err)
		}
		fmt.Println("recv:", recv.Payload)
		err = ss.SendMsg(recv)
		if err != nil {
			panic(err)
		}
	}
	return err
}

func (q QuicServer) Close() error {
	//TODO implement me
	panic("implement me")
}

func main() {
	starter := api.NewServer()
	starter.Register(newQuicServer, nil)
	starter.Serve()
}
