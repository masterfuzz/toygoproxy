.PHONY: build docker-build generate test

NAME := toygoproxy
VERSION ?= dev


clean:
	rm -r bin/ || true

build:
	go build -o bin/$(NAME) cmd/main.go

docker-build: clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(NAME) cmd/main.go
	chmod +x bin/$(NAME)
	docker build --load --platform linux/amd64 -t $(NAME):$(VERSION) . --no-cache

# docker-push:
# 	docker push $(HUB)/$(NAME):$(VERSION)

generate:
	sqlc generate -f pkg/database/sqlc.yaml

test:
	go test -v ./...

format:
	go fmt ./...
