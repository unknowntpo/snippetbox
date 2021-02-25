BIN=web
GOROOT=$(shell go env GOROOT)

all: build

build:
	go build -o $(BIN) ./cmd/web
tls:
	@echo "Generating a Self-Signed TLS Certificate"
	cd ./tls && go run $(GOROOT)/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

clean:
	-rm -r $(BIN) ./tls/*.pem

.PHONY: build tls clean

