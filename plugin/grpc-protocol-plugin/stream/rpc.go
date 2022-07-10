package stream

import (
	"bytes"
	"encoding/binary"
	"github.com/zhengheng7913/grpc-go-starter/pkg/codec"
	"io"
	"math"
)

const (
	sizeLen   = 4
	headerLen = sizeLen
)

type composer struct {
	w io.Writer
}

func (c *composer) sendMsg(b []byte) error {
	_, err := c.w.Write(b)
	return err
}

type parser struct {
	r io.Reader
	// 读取当前数据报的长度
	header [headerLen]byte
}

func (p *parser) recvMsg(maxRecvLength uint32) ([]byte, error) {
	if _, err := p.r.Read(p.header[:]); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(p.header[:])
	if length == 0 {
		return nil, nil
	}
	if length > math.MaxUint32 {
		return nil, io.EOF
	}
	if length > maxRecvLength {
		return nil, io.EOF
	}
	msg := make([]byte, int(length))
	if _, err := p.r.Read(msg); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	return msg, nil
}

func rpcRecv(m interface{}, p *parser, c codec.Codec, maxRecvLength uint32) error {
	b, err := p.recvMsg(maxRecvLength)
	if err != nil {
		return err
	}
	if err = c.Unmarshal(b, m); err != nil {
		return err
	}
	return nil
}

func rpcSend(m interface{}, c *composer, cc codec.Codec, maxSendLength uint32) (err error) {
	data, err := cc.Marshal(m)
	if err != nil {
		return err
	}
	if uint32(len(data)) > maxSendLength {
		return io.EOF
	}

	dl := uint32(len(data))
	hdr := make([]byte, headerLen)
	binary.BigEndian.PutUint32(hdr[:], dl)
	buf := bytes.Buffer{}
	buf.Write(hdr)
	buf.Write(data)
	return c.sendMsg(buf.Bytes())
}
