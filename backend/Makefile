PROJECT_NAME=mmbrnr

MOCKS_DIR=./test/mocks
MOCKS_GEN_ARGS=--case underscore --keeptree --all --exported

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

.PHONY: gen-mocks
gen-mocks:
	mokcery $(MOCKS_GEN_ARGS)
	

.PHONY: unit-test
unit-test: 
	go test -cover -coverprofile test/coverage.out -count=1 $(FILE_GLOB)
	go tool cover -o=test/coverage.html -html=test/coverage.out