
.PHONY: all
all: fmt test lint build

.PHONY: lint
lint:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	CGO_ENABLED=0 go build -o build/paranotify main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/paranotify_amd64 main.go

.PHONY: server
server:
	export $$(cat .env) && go run main.go -token=$${TOKEN} -reply-server

.PHONY: echo
echo:
	export $$(cat .env) && echo "this is echo msg" | go run main.go -token=$${TOKEN} -user-id=$${USER}

.PHONY: install
install: build
	sudo cp build/paranotify /usr/local/bin/paranotify

.PHONY: docker
docker:
	docker build . -t paragor/paranotify:latest
	docker push paragor/paranotify:latest
