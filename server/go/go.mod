module github.com/s77rt/rdpcloud/server/go

go 1.19

replace github.com/s77rt/rdpcloud/proto/go => ../../proto/go

require (
	github.com/s77rt/rdpcloud/proto/go v0.0.0-00010101000000-000000000000
	golang.org/x/sys v0.0.0-20220906165534-d0df966e6959
	google.golang.org/grpc v1.49.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
