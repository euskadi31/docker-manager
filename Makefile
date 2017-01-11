.PHONY: all clean deps fmt vet test docker

EXECUTABLE ?= docker-manager
IMAGE ?= euskadi31/$(EXECUTABLE)
VERSION ?= $(shell git describe --match 'v[0-9]*' --dirty='-dev' --always)
COMMIT ?= $(shell git rev-parse --short HEAD)

LDFLAGS = -X "main.Revision=$(COMMIT)" -X "main.Version=$(VERSION)"
PACKAGES = $(shell go list ./... | grep -v /vendor/)

release:
	@echo "Release v$(version)"
	@git pull
	@git checkout master
	@git pull
	@git checkout develop
	@git flow release start $(version)
	@echo "$(version)" > .version
	@git add .version
	@git commit -m "feat(project): update version file" .version
	@git flow release finish $(version) -p -m "Release v$(version)"
	@git checkout develop
	@echo "Release v$(version) finished."

all: deps build test

clean:
	@go clean -i ./...

deps:
	@glide install

fmt:
	@go fmt $(PACKAGES)

vet:
	@go vet $(PACKAGES)

test:
	@for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

docker:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w $(LDFLAGS)'
	@docker build --rm -t $(IMAGE) .

$(EXECUTABLE): $(wildcard *.go)
	@go build -ldflags '-s -w $(LDFLAGS)'

build: $(EXECUTABLE)

run: build
	@PORT=1339 DEBUG=true ./$(EXECUTABLE)
