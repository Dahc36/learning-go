.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

vet: fmt
	go vet ./...
.PHONY:vet

build: vet
	go build -o build.exe hello.go
.PHONY:build

run: vet
	go run hello.go
.PHONY:run
