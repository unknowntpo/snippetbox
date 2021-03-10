BIN=web
GOROOT=$(shell go env GOROOT)

all: 

run-prod:
	$ sudo docker-compose -f docker-compose-prod.yaml up -d --build
stop-prod:
	$ sudo docker-compose down
run-dev:
	$ sudo docker-compose up -d --force-recreate --build --always-recreate-deps
stop-dev:
	$ sudo docker-compose down -v

build: 
	go build -o $(BIN) ./cmd/web
tls-prod:
	@echo "Use production TLS Certificate"
	-rm -rf ./tls
	-mv ./tls-prod ./tls
tls-dev:
	@echo "Generating a Self-Signed TLS Certificate"
	-rm -rf ./tls
	-mkdir ./tls && cd ./tls && go run $(GOROOT)/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

clean:
	-rm -r $(BIN) ./tls

.PHONY: build tls-prod tls-dev clean

