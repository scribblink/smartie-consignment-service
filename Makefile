build: 
	protoc -I. --go_out=plugins=micro:. \
		proto/consignment/consignment.proto
	GOOS=linux GOARCH=amd64 go build -o smartie-consignment-service && \
	docker build -t smartie-consignment-service .

run:
	docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 smartie-consignment-service