GIT_TAG = "$(shell git describe --tags --always)"

export SERVER_NAME
export SERVER_IP
export IS_FREE_TRIAL
export FREE_TRIAL_DURATION = 30

ifndef SERVER_NAME
$(error SERVER_NAME is not set)
endif
ifndef SERVER_IP
$(error SERVER_IP is not set)
endif
ifndef IS_FREE_TRIAL
$(error IS_FREE_TRIAL is not set)
endif

ifeq ($(IS_FREE_TRIAL), FALSE)
	LICENSE_EXP_DATE = ""
else
ifeq ($(IS_FREE_TRIAL), TRUE)
	LICENSE_EXP_DATE = "$(shell date -d '+$(FREE_TRIAL_DURATION) day' '+%Y-%m-%d')"
else
$(error IS_FREE_TRIAL can either be TRUE or FALSE)
endif
endif

SERVER_LD_FLAGS = "-X 'main.Version=$(GIT_TAG)' -X 'main.ServerName=$(SERVER_NAME)' -X 'main.ServerIP=$(SERVER_IP)' -X 'main.LicenseExpDate=$(LICENSE_EXP_DATE)'"
CLIENT_LD_FLAGS = "-X 'main.Version=$(GIT_TAG)' -X 'main.ServerName=$(SERVER_NAME)' -X 'main.ServerIP=$(SERVER_IP)' -X 'main.LicenseExpDate=$(LICENSE_EXP_DATE)'"

all: info gen-cert gen-go gen-php build-server-go build-client-frontend-react build-client-go build-client-php-whmcs-provisioning-module info

info:
	@echo "SERVER_NAME: $(SERVER_NAME)"
	@echo "SERVER_IP: $(SERVER_IP)"
	@echo "IS_FREE_TRIAL: $(IS_FREE_TRIAL) ($(FREE_TRIAL_DURATION) days)"
	@echo "LICENSE_EXP_DATE: $(LICENSE_EXP_DATE)"

gen-cert:
	cd cert && bash gen.sh

gen-go:
	rm -rf proto/go && mkdir -p proto/go && touch proto/go/.keep
	cd proto/go && go mod init github.com/s77rt/rdpcloud/proto/go
	protoc --proto_path=proto/proto --go_out=proto/go --go_opt=paths=source_relative --go-grpc_out=proto/go --go-grpc_opt=paths=source_relative proto/proto/*/*/*.proto

gen-php:
	rm -rf proto/php && mkdir -p proto/php && touch proto/php/.keep
	protoc --proto_path=proto/proto --php_out=proto/php --grpc_out=proto/php --plugin=protoc-gen-grpc=/usr/local/bin/grpc_php_plugin proto/proto/*/*/*.proto

build-server-go:
	rm -rf server/go/bin && mkdir -p server/go/bin && touch server/go/bin/.keep
	rm -rf server/go/cert && mkdir -p server/go/cert && touch server/go/cert/.keep
	cp cert/server-cert.pem server/go/cert/
	cp cert/server-key.pem server/go/cert/
	cd server/go && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags $(SERVER_LD_FLAGS) -o bin/rdpcloud-server-windows.exe
	cd server/go/bin && 7za -y a rdpcloud-server-windows.7z rdpcloud-server-windows.exe

build-client-frontend-react:
	npm --prefix client-frontend/react/rdpcloud-client install
	npm --prefix client-frontend/react/rdpcloud-client run build

build-client-go:
	rm -rf client/go/bin && mkdir -p client/go/bin && touch client/go/bin/.keep
	rm -rf client/go/cert && mkdir -p client/go/cert && touch client/go/cert/.keep
	cp cert/server-cert.pem client/go/cert/
	cp client-frontend/react/rdpcloud-client/build/static/js/main.*.js client/go/frontend/static/assets/js/app/main.js
	cp client-frontend/react/rdpcloud-client/build/static/js/*.chunk.js client/go/frontend/static/assets/js/app/chunk.js
	cp client-frontend/react/rdpcloud-client/build/static/css/main.*.css client/go/frontend/static/assets/css/app/main.css
	cd client/go && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags $(CLIENT_LD_FLAGS) -o bin/rdpcloud-client-linux
	cd client/go/bin && 7za -y a rdpcloud-client-linux.7z rdpcloud-client-linux

build-client-php-whmcs-provisioning-module:
	# src
	rm -rf client/php/whmcs-provisioning-module/src/rdpcloud/lib/proto && mkdir -p client/php/whmcs-provisioning-module/src/rdpcloud/lib/proto && touch client/php/whmcs-provisioning-module/src/rdpcloud/lib/proto/.keep
	cp -r proto/php/* client/php/whmcs-provisioning-module/src/rdpcloud/lib/proto
	cd client/php/whmcs-provisioning-module && bash init_composer.sh
	rm -rf client/php/whmcs-provisioning-module/src/rdpcloud/vendor && mkdir -p client/php/whmcs-provisioning-module/src/rdpcloud/vendor && touch client/php/whmcs-provisioning-module/src/rdpcloud/vendor/.keep
	rm -rf client/php/whmcs-provisioning-module/src/rdpcloud/cert && mkdir -p client/php/whmcs-provisioning-module/src/rdpcloud/cert && touch client/php/whmcs-provisioning-module/src/rdpcloud/cert/.keep
	cp cert/server-cert.pem client/php/whmcs-provisioning-module/src/rdpcloud/cert/
	# dist
	rm -rf client/php/whmcs-provisioning-module/dist && mkdir -p client/php/whmcs-provisioning-module/dist && touch client/php/whmcs-provisioning-module/dist/.keep
	cd client/php/whmcs-provisioning-module/src/rdpcloud && composer install --no-dev --optimize-autoloader
	cd client/php/whmcs-provisioning-module && mkdir -p dist/rdpcloud
	cd client/php/whmcs-provisioning-module && cp README.md dist/rdpcloud/
	cd client/php/whmcs-provisioning-module && cp src/htaccess/.htaccess dist/rdpcloud/
	cd client/php/whmcs-provisioning-module && cp -r src/rdpcloud/lib/ dist/rdpcloud/
	cd client/php/whmcs-provisioning-module && cp -r src/rdpcloud/vendor/ dist/rdpcloud/
	cd client/php/whmcs-provisioning-module && cp -r src/rdpcloud/cert/ dist/rdpcloud/
	cd client/php/whmcs-provisioning-module && rm dist/rdpcloud/lib/proto/.keep
	cd client/php/whmcs-provisioning-module && rm dist/rdpcloud/vendor/.keep
	cd client/php/whmcs-provisioning-module && rm dist/rdpcloud/cert/.keep
	cd client/php/whmcs-provisioning-module && cp src/rdpcloud/rdpcloud.php dist/rdpcloud/
	cd client/php/whmcs-provisioning-module/dist && zip -r rdpcloud.zip rdpcloud
