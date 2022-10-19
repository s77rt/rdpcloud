GIT_TAG := $(shell git describe --tags --always)

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
	EXP_DATE :=
else
ifeq ($(IS_FREE_TRIAL), TRUE)
	EXP_DATE := $(shell date -d '+$(FREE_TRIAL_DURATION) day' '+%Y-%m-%d')
else
$(error IS_FREE_TRIAL can either be TRUE or FALSE)
endif
endif

ENCRYPTION_KEY := $(shell openssl rand 128 | base64 -w 0)
ENCRYPTION_KEY_X :=  $(shell xor '$(ENCRYPTION_KEY)' 'RDPCloud' | base64 -w 0)
ENCRYPTED_SERVER_NAME := $(shell xor '$(SERVER_NAME)' '$(ENCRYPTION_KEY_X)' | base64 -w 0)
ENCRYPTED_SERVER_IP := $(shell xor '$(SERVER_IP)' '$(ENCRYPTION_KEY_X)' | base64 -w 0)
ENCRYPTED_EXP_DATE := $(shell xor '$(EXP_DATE)' '$(ENCRYPTION_KEY_X)' | base64 -w 0)
SIGNATURE := $(shell xor '$(ENCRYPTED_SERVER_NAME)$(ENCRYPTED_SERVER_IP)$(ENCRYPTED_EXP_DATE)SIGNATURE' '$(ENCRYPTION_KEY_X)' | base64 -w 0)

SERVER_GO_LDFLAGS := -X 'main.Version=$(GIT_TAG)' -X 'github.com/s77rt/rdpcloud/server/go/license.EncryptionKey=$(ENCRYPTION_KEY)' -X 'github.com/s77rt/rdpcloud/server/go/license.EncryptedServerName=$(ENCRYPTED_SERVER_NAME)' -X 'github.com/s77rt/rdpcloud/server/go/license.EncryptedServerIP=$(ENCRYPTED_SERVER_IP)' -X 'github.com/s77rt/rdpcloud/server/go/license.EncryptedExpDate=$(ENCRYPTED_EXP_DATE)' -X 'github.com/s77rt/rdpcloud/server/go/license.Signature=$(SIGNATURE)'
CLIENT_GO_LDFLAGS := -X 'main.Version=$(GIT_TAG)' -X 'github.com/s77rt/rdpcloud/client/go/license.EncryptionKey=$(ENCRYPTION_KEY)' -X 'github.com/s77rt/rdpcloud/client/go/license.EncryptedServerName=$(ENCRYPTED_SERVER_NAME)' -X 'github.com/s77rt/rdpcloud/client/go/license.EncryptedServerIP=$(ENCRYPTED_SERVER_IP)' -X 'github.com/s77rt/rdpcloud/client/go/license.EncryptedExpDate=$(ENCRYPTED_EXP_DATE)' -X 'github.com/s77rt/rdpcloud/client/go/license.Signature=$(SIGNATURE)'

all: dep gen-cert gen-go gen-php build-server-go build-client-frontend-react build-client-go build-client-php-whmcs-provisioning-module build-bundle info

dep:
	go install mvdan.cc/garble@latest
	go install github.com/s77rt/xor/cmd/xor@latest

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
	cd server/go && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 garble -literals -tiny -seed=random build -ldflags '$(SERVER_GO_LDFLAGS)' -o bin/rdpcloud-server-windows.exe
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
	cd client/go && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 garble -literals -tiny -seed=random build -ldflags '$(CLIENT_GO_LDFLAGS)' -o bin/rdpcloud-client-linux
	cd client/go/bin && 7za -y a rdpcloud-client-linux.7z rdpcloud-client-linux
	cd client/go && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 garble -literals -tiny -seed=random build -ldflags '$(CLIENT_GO_LDFLAGS)' -o bin/rdpcloud-client-windows.exe
	cd client/go/bin && 7za -y a rdpcloud-client-windows.7z rdpcloud-client-windows.exe
	cd client/go && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 garble -literals -tiny -seed=random build -ldflags '$(CLIENT_GO_LDFLAGS)' -o bin/rdpcloud-client-darwin
	cd client/go/bin && 7za -y a rdpcloud-client-darwin.7z rdpcloud-client-darwin

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
	cd client/php/whmcs-provisioning-module && cp src/htaccess/.htaccess dist/rdpcloud/
	cd client/php/whmcs-provisioning-module && cp -r src/rdpcloud/lib/ dist/rdpcloud/
	cd client/php/whmcs-provisioning-module && cp -r src/rdpcloud/vendor/ dist/rdpcloud/
	cd client/php/whmcs-provisioning-module && cp -r src/rdpcloud/cert/ dist/rdpcloud/
	cd client/php/whmcs-provisioning-module && cp -r src/rdpcloud/templates/ dist/rdpcloud/
	cd client/php/whmcs-provisioning-module && rm dist/rdpcloud/lib/proto/.keep
	cd client/php/whmcs-provisioning-module && rm dist/rdpcloud/vendor/.keep
	cd client/php/whmcs-provisioning-module && rm dist/rdpcloud/cert/.keep
	cd client/php/whmcs-provisioning-module && cp src/rdpcloud/rdpcloud.php dist/rdpcloud/
	cd client/php/whmcs-provisioning-module/dist && zip -r rdpcloud.zip rdpcloud

build-bundle:
	rm -rf bundle && mkdir bundle && touch bundle/.keep
	cp docs/EULA.md bundle/
	cp -r docs/ bundle/
	mkdir bundle/server
	cp server/go/bin/*.7z bundle/server/
	cp server/go/README.md bundle/server/
	mkdir bundle/client
	cp client/go/bin/*.7z bundle/client/
	cp client/go/README.md bundle/client/
	mkdir -p bundle/integration/whmcs-provisioning-module
	cp client/php/whmcs-provisioning-module/dist/*.zip  bundle/integration/whmcs-provisioning-module/
	cp client/php/whmcs-provisioning-module/README.md bundle/integration/whmcs-provisioning-module/
	mkdir bundle/development
	cp -r proto/ bundle/development/
	mkdir bundle/development/cert
	cp cert/server-cert.pem bundle/development/cert/
	cd bundle && echo "${SERVER_NAME} (${SERVER_IP})" > LICENSEE.md
	cd bundle && echo "Read the docs inside the /docs folder" > README.md
	cd bundle && 7za -y a -x!.keep rdpcloud-bundle.7z .

info:
	@echo '########################################'
	@echo 'SERVER_NAME:     $(SERVER_NAME)'
	@echo 'SERVER_IP:       $(SERVER_IP)'
	@echo 'IS_FREE_TRIAL:   $(IS_FREE_TRIAL) ($(FREE_TRIAL_DURATION) days)'
	@echo 'EXP_DATE:        $(EXP_DATE)'
	@echo '########################################'
