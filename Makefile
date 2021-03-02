BIN=web
GOROOT=$(shell go env GOROOT)

all: build tls

compose: 
	$ docker-compose up --force-recreate
build: 
	go build -o $(BIN) ./cmd/web
tls:
	@echo "Generating a Self-Signed TLS Certificate"
	-rm -rf ./tls
	-mkdir ./tls && cd ./tls && go run $(GOROOT)/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

clean:
	-rm -r $(BIN) ./tls

.PHONY: build tls clean

