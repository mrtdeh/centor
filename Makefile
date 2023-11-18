

.PHONY: all test clean

all: test build

test:
	go test ./...

build: clean
	go build -v -o bin/centor  ./main.go

protoc:
	protoc --go-grpc_out=require_unimplemented_servers=false:./proto/ ./proto/*.proto --go_out=./proto


build-compose:
	docker-compose up -d --build

clean:
	rm -f ./bin/*