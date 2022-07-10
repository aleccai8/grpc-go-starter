package grpc_protocol_plugin

import (
	"encoding/json"
	"fmt"
	"github.com/zhengheng7913/grpc-go-starter/pkg/codec"
	"google.golang.org/protobuf/proto"
)

func init() {
	codec.Register("json", &JSONCodec{}, &JSONCodec{})
	codec.Register("proto", &ProtoCodec{}, &ProtoCodec{})
}

type JSONCodec struct {
}

func (c *JSONCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type ProtoCodec struct {
}

func (c *ProtoCodec) Marshal(v interface{}) ([]byte, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}
	return proto.Marshal(vv)
}

func (c *ProtoCodec) Unmarshal(data []byte, v interface{}) error {
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}
	return proto.Unmarshal(data, vv)
}
