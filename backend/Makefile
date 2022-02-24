PROJECT_NAME=mmbrnr

ifndef FILE_GLOB
	FILE_GLOB=./...
endif

.PHONY: build
build:
	go mod tidy
	go build -o /go/bin/$(PROJECT_NAME) main.go

.PHONY: serve
serve: build
	$(PROJECT_NAME) serve

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: unit-test
unit-test: 
	go test -cover -count=1 $(FILE_GLOB)