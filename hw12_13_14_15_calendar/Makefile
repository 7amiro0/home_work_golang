BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	mkdir -p ./google/api
	(test -f ./google/api/annotations.proto) || curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > "./google/api/annotations.proto"
	(test -f ./google/api/http.proto) || curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > "./google/api/http.proto"


	protoc -I=. \
		--go_out ./internal/server/grpc --go_opt=paths=source_relative \
		--go-grpc_out ./internal/server/grpc --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=./internal/server/grpc/ \
		--grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		./google/EventService.proto

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) --config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	 run ./...


.PHONY: build run build-img run-img version test lint
