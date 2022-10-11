module github.com/s77rt/rdpcloud/client/go

go 1.19

replace github.com/s77rt/rdpcloud/proto/go => ../../proto/go

require (
	github.com/fullstorydev/grpcui v1.3.1
	github.com/s77rt/xor v0.0.0-20221010224322-0f0d4971e11f
	google.golang.org/grpc v1.49.0
)

require (
	github.com/fullstorydev/grpcurl v1.8.6 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jhump/protoreflect v1.12.0 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20220406163625-3f8b81556e12 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
