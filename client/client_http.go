package client

import "context"

const HttpProtocol = "http"

func NewHttpClient() Client {
	return &HttpClient{}
}

type HttpClient struct {
}

func (h *HttpClient) Invoke(context context.Context, method string, req interface{}, opts ...Option) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HttpClient) Register(realClient interface{}, opts ...Option) {
	//TODO implement me
	panic("implement me")
}
