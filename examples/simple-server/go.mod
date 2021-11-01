module github.com/zhengheng7913/grpc-config/examples/simple-server

go 1.16

require (
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.6.0
	github.com/polarismesh/polaris-go v1.0.0 // indirect
	github.com/zhengheng7913/grpc-config v0.0.0
	github.com/zhengheng7913/grpc-polaris-plugin v0.0.0
	google.golang.org/genproto v0.0.0-20210903162649-d08c68adba83
)

replace github.com/zhengheng7913/grpc-config v0.0.0 => ../../../grpc-config

replace github.com/zhengheng7913/grpc-polaris-plugin v0.0.0 => /Volumes/Develop/grpc-polaris-plugin
