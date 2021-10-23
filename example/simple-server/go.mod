module github.com/zhengheng7913/grpc-config/example/simple-server

go 1.17

require (
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210903162649-d08c68adba83 // indirect
)

require grpc-config v0.0.0

replace grpc-config v0.0.0 => ../../../grpc-config
