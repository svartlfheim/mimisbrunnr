PROJECT_NAME=mmbrnr

.PHONY: build
build:
	go mod tidy
	go build -o /go/bin/$(PROJECT_NAME) main.go

.PHONY: serve
serve: build
	$(PROJECT_NAME) serve