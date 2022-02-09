package client

import "context"

const HttpProtocol = "http"

func NewHttpClient() Client {
	return &HttpClient{}
}

type HttpClient struct {
}

func (h HttpClient) Invoke(context context.Context, method any, req any, options *Options) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (h HttpClient) RealClient() any {
	//TODO implement me
	panic("implement me")
}

func (h HttpClient) Register(realClient any, options *Options) {
	//TODO implement me
	panic("implement me")
}
