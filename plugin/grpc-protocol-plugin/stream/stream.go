package stream

import (
	"github.com/zhengheng7913/grpc-go-starter/pkg/codec"
	"io"
)

type Stream interface {
	RecvMsg(m interface{}) error
	SendMsg(m interface{}) error
}

func NewStream(r io.Reader, w io.Writer, cc codec.Codec) Stream {
	return &clientStream{
		c:             &composer{w: w},
		p:             &parser{r: r},
		cc:            cc,
		maxRecvLength: 1024,
		maxSendLength: 1024,
	}
}

type clientStream struct {
	cc            codec.Codec
	c             *composer
	p             *parser
	maxRecvLength uint32
	maxSendLength uint32
}

func (c *clientStream) RecvMsg(m interface{}) error {
	return rpcRecv(m, c.p, c.cc, c.maxRecvLength)
}

func (c *clientStream) SendMsg(m interface{}) error {
	return rpcSend(m, c.c, c.cc, c.maxSendLength)
}
