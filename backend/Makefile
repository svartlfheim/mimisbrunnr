PROJECT_NAME=mmbrnr

MOCKS_DIR=./test/mocks
MOCKS_GEN_ARGS=--case underscore --keeptree --all --with-expecter --exported

ifeq ($(V),1)
	TEST_ARGS="-v"
endif

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

.PHONY: gen-cmd-mocks
gen-cmd-mocks: 
	mockery $(MOCKS_GEN_ARGS) --output $(MOCKS_DIR)/cmd --outpkg cmdmocks --dir ./cmd

.PHONY: gen-internal-schema-mocks
gen-internal-schema-mocks: 
	mockery --case underscore --with-expecter --exported --output $(MOCKS_DIR)/schema --outpkg schema --dir ./internal/schema --name connectionManager
	mockery --case underscore --with-expecter --exported --output $(MOCKS_DIR)/schema --outpkg schema --dir ./internal/schema --name hasSchema

.PHONY: gen-scm-postgres-mocks
gen-scm-postgres-mocks: 
	mockery --case underscore --with-expecter --exported --output $(MOCKS_DIR)/scm/postgres --outpkg scmpostgresmocks --dir ./internal/scm/postgres --name connManager

.PHONY: gen-scm-mocks
gen-scm-mocks: gen-scm-postgres-mocks

.PHONY: gen-internal-rdbconn-mocks
gen-internal-rdbconn-mocks: 
	mockery --case underscore --with-expecter --exported --output $(MOCKS_DIR)/rdbconn --outpkg rdbconnmocks --dir ./internal/rdbconn --name connectionOpener

.PHONY: gen-internal-validation-mocks
gen-internal-validation-mocks: 
	mockery --case underscore --with-expecter --exported --output $(MOCKS_DIR)/validation --outpkg validationmocks --dir ./internal/validation --name ValidationError

.PHONY: gen-mocks
gen-mocks: gen-cmd-mocks gen-internal-rdbconn-mocks gen-internal-validation-mocks gen-internal-schema-mocks gen-scm-mocks

.PHONY: unit-test
unit-test: 
	go test -cover -coverprofile test/unit.out $(TEST_ARGS) $(FILE_GLOB)
	go tool cover -o=test/unit.html -html=test/unit.out

.PHONY: integration-test
integration-test:
	CI_INTEGRATION_TESTS_ENABLED=true go test -cover -coverprofile test/integration.out -p 1 -run "^TestIntegration" $(TEST_ARGS) $(FILE_GLOB)
	go tool cover -o=test/integration.html -html=test/integration.out

.PHONY: test
test: unit-test integration-test

.PHONY: ci
ci: fmt lint test