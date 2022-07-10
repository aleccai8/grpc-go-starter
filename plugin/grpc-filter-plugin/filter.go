package filter

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	starter "github.com/zhengheng7913/grpc-go-starter/pkg/filter"
	"github.com/zhengheng7913/grpc-go-starter/plugin/grpc-filter-plugin/error"
	"github.com/zhengheng7913/grpc-go-starter/plugin/grpc-filter-plugin/http"
	"github.com/zhengheng7913/grpc-go-starter/plugin/grpc-filter-plugin/validate"
	"google.golang.org/grpc"
)

func init() {
	starter.Register(error.HttpErrorFilter, runtime.WithErrorHandler(error.HttpErrorHandler), nil)
	starter.Register(http.KeyLocation302, runtime.WithForwardResponseOption(http.Location302), nil)
	starter.Register(error.GrpcRecoverServerFilter, grpc.ChainUnaryInterceptor(error.UnaryServerInterceptor), nil)
	starter.Register(validate.GrpcUnaryServerFilter, grpc.ChainUnaryInterceptor(validate.UnaryServerInterceptor), nil)
}
