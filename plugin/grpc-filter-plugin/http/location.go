package http

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
	"net/http"
)

const (
	KeyLocation302     = "http-location-302"
	KeyGrpcLocation302 = "Grpc-Metadata-Http-Location-302"
)

func Location302(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set http status code
	if vals := md.HeaderMD.Get(KeyLocation302); len(vals) > 0 {
		w.Header().Set("location", vals[0])
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, KeyLocation302)
		delete(w.Header(), KeyGrpcLocation302)
		w.WriteHeader(302)
	}

	return nil
}
