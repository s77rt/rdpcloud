module github.com/s77rt/rdpcloud/server/go

go 1.19

replace github.com/s77rt/rdpcloud/proto/go => ../../proto/go

require (
	github.com/beevik/ntp v0.3.0
	github.com/go-co-op/gocron v1.17.0
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/s77rt/rdpcloud/proto/go v0.0.0-00010101000000-000000000000
	github.com/s77rt/xor v0.0.0-20221010224322-0f0d4971e11f
	golang.org/x/sys v0.0.0-20220906165534-d0df966e6959
	google.golang.org/grpc v1.49.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
