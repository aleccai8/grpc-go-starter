package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zhengheng7913/grpc-go-starter/filter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	HeaderErrorMsg  = "grpc-error-msg"
	HeaderErrorCode = "grpc-error-code"
	HttpErrorFilter = "http-error"
)

func init() {
	filter.Register(HttpErrorFilter, runtime.WithErrorHandler(httpErrorHandler), nil)
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
