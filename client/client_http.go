package client

import "context"

const HttpProtocol = "http"

func NewHttpClient() Client {
	return &HttpClient{}
}

type HttpClient struct {
}

func (h HttpClient) RealClient(ctx context.Context) any {
	//TODO implement me
	panic("implement me")
}

func (h HttpClient) Register(realClient any, opts ...Option) {
	//TODO implement me
	panic("implement me")
}
