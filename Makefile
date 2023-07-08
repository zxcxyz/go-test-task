install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

lint:
	golangci-lint run .

find:
	find / -name golangci-lint 2>/dev/null

build:
	docker build -f Dockerfile -t appserv .

dockerrun:
	docker run -p:8080:8080 appserv

make update:
	docker build -f Dockerfile -t appserv .
	docker run -p:8080:8080 appserv

run:
	SERVER_PORT=8080 METRICS_ENDPOINT=/metrics go run main.go