package error

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zhengheng7913/grpc-config/filter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

const (
	HeaderErrorMsg  = "grpc-error-msg"
	HeaderErrorCode = "grpc-error-code"
)

type Handler struct {
}

func (h Handler) Apply(arg ...interface{}) {
	panic("implement me")
}

func (h Handler) ServiceName() string {
	return filter.GlobalFilterName
}

func (h Handler) ProtocolName() string {
	return "http"
}

func httpErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {

	var customStatus *runtime.HTTPStatusError
	if errors.As(err, &customStatus) {
		err = customStatus.Err
	}

	s := status.Convert(err)

	w.Header().Del("Trailer")
	w.Header().Del("Transfer-Encoding")

	if s.Code() == codes.Unauthenticated {
		w.Header().Set("WWW-Authenticate", s.Message())
	}

	w.Header().Set(HeaderErrorMsg, s.Message())
	w.Header().Set(HeaderErrorCode, strconv.Itoa(int(s.Code())))

	st := runtime.HTTPStatusFromCode(s.Code())
	if customStatus != nil {
		st = customStatus.HTTPStatus
	}

	w.WriteHeader(st)
}
