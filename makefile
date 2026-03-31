.PHONY: help run test test-v fmt vet tidy up down logs rebuild

help:
	@echo "Targets:"
	@echo " make run		- run app locally"
	@echo " make test		- run uint-tests"
	@echo " make test-v		- run tests with verbose output"
	@echo " make fmt		- format go code"
	@echo " make vet		- run go vet"
	@echo " make tidy		- tidy go modules"
	@echo " make up			- docker-compose up -d"
	@echo " make down		- docker-compose down"
	@echo " make logs		- show docker logs"
	@echo " make rebuild	- rebuild and restart containers"

run:
	go run ./cmd

test:
	go test ./...

test-v:
	go test ./... -v

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f --tail=100

rebuild:
	docker-compose down
	docker-compose up -d --build