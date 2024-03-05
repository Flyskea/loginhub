version:=$(shell git describe --tags --always HEAD)
flags=-X main.Version=$(version)

.PHONY: serverdocker
serverdocker:
	docker build -f ./deploy/docker/Dockerfile -t loginhub:$(version) --build-arg FLAGS="$(flags)" .

.PHONY: wire
wire:
	cd cmd/server && wire

.PHONY: build
build:
	go build -ldflags "$(flags)" -o bin/loginhub ./cmd/server

.PHONY: doc
doc: 
	swag fmt
	swag init -g ./internal/server/http/apiv1.go -dir ./ --parseInternal
