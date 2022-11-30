CALENDAR := "./bin/calendar"
SCHEDULER := "./bin/scheduler"
SENDER := "./bin/sender"
DOCKER_IMG="calendar:develop"

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
	go build -v -o $(CALENDAR) ./cmd/calendar
	go build -v -o $(SCHEDULER) ./cmd/scheduler
	go build -v -o $(SENDER) ./cmd/sender

run: build
	$(CALENDAR) --config ./configs/calendar_config.yaml &
	$(SCHEDULER) --config ./configs/scheduler_config.yaml &
	$(SENDER) --config ./configs/sender_config.yaml

db:
	systemctl start postgresql.service

docker:
	systemctl unmask docker && systemctl start docker && systemctl status docker

rabbit:
	sudo docker rm -f rabbitmq
	sudo docker run -d --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management

up:
	sudo docker-compose -f ./deployments/docker-compose.yaml -f ./deployments/docker-compose.prod.yaml up --build

down:
	sudo docker-compose -f ./deployments/docker-compose.yaml down

version: build
	$(CALENDAR) version

test:
	go test -race ./internal/...

integration-tests:
	sudo docker-compose -f ./deployments/docker-compose.yaml -f ./deployments/docker-compose.tests.yaml up --build

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	 run ./...


.PHONY: build run build-img run-img version test lint
