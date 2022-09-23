GIT_TAG = "$(shell git describe --tags --always)"
DISPLAY_NAME ?= RDPCloud

SERVER_LD_FLAGS = "-X 'main.Version=$(GIT_TAG)'"
CLIENT_LD_FLAGS = "-X 'main.Version=$(GIT_TAG)' -X 'main.DisplayName=$(DISPLAY_NAME)'"

all: gen-go build-server-go-windows build-client-frontend-react link-client-frontend-react-with-go build-client-go-linux

gen-go:
	rm -rf proto/go && mkdir proto/go
	cd proto/go && go mod init github.com/s77rt/rdpcloud/proto/go
	protoc --proto_path=proto/proto --go_out=proto/go --go_opt=paths=source_relative --go-grpc_out=proto/go --go-grpc_opt=paths=source_relative proto/proto/*/*/*.proto

build-server-go-windows:
	rm -rf server/go/bin/windows && mkdir -p server/go/bin/windows
	cd server/go && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags $(SERVER_LD_FLAGS) -o bin/windows/rdpcloud_server.exe
	cd server/go/bin/windows && 7za -y a rdpcloud_server_windows.7z rdpcloud_server.exe

build-client-go-linux:
	rm -rf client/go/bin/linux && mkdir -p client/go/bin/linux
	cd client/go && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags $(CLIENT_LD_FLAGS) -o bin/linux/rdpcloud_client
	cd client/go/bin/linux && 7za -y a rdpcloud_client_linux.7z rdpcloud_client

build-client-frontend-react:
	npm --prefix client-frontend/react/rdpcloud-client install
	npm --prefix client-frontend/react/rdpcloud-client run build

link-client-frontend-react-with-go:
	cp client-frontend/react/rdpcloud-client/build/static/js/main.*.js client/go/frontend/static/assets/js/app/main.js
	cp client-frontend/react/rdpcloud-client/build/static/js/*.chunk.js client/go/frontend/static/assets/js/app/chunk.js
	cp client-frontend/react/rdpcloud-client/build/static/css/main.*.css client/go/frontend/static/assets/css/app/main.css
