

.PHONY: all test clean

all: test build

test:
	go test ./...

build: clean
	go build -v -o bin/centor  ./main.go

protoc:
	protoc --go-grpc_out=require_unimplemented_servers=false:./proto/ ./proto/*.proto --go_out=./proto



docker-clean:
	docker image prune -f
	
docker-build: docker-clean
	docker build . --tag mrtdeh/centor
docker-up:
	docker compose -p dc1 up --force-recreate --build -d

docker-up-dc2: 
	docker compose -p dc2 -f ./docker-compose-dc2.yml up --force-recreate -d

	


clean:
	rm -f ./bin/*