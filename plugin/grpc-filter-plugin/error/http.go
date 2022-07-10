package error

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

const (
	HeaderErrorMsg  = "grpc-error-msg"
	HeaderErrorRet  = "grpc-error-ret"
	HttpErrorFilter = "http-error"
)

func init() {

}

func HttpErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	var customStatus *runtime.HTTPStatusError
	fmt.Print(err)
	if errors.As(err, &customStatus) {
		err = customStatus.Err
	}

	s := status.Convert(err)

	if s.Code() == codes.Unauthenticated {
		w.Header().Set("WWW-Authenticate", s.Message())
	}

	w.Header().Set(HeaderErrorMsg, s.Message())
	w.Header().Set(HeaderErrorRet, strconv.Itoa(int(s.Code())))

	if s.Code() < 1000 {
		w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	} else {
		w.WriteHeader(200)
	}

}
