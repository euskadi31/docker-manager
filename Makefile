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
	@sudo docker build --rm -t $(IMAGE) .

publish: docker
	@sudo docker tag $(IMAGE) $(IMAGE):latest
	@sudo docker push $(IMAGE)

$(EXECUTABLE): $(wildcard *.go)
	@echo "Building $(EXECUTABLE)..."
	@go build -ldflags '-s -w $(LDFLAGS)'

build: $(EXECUTABLE)

run: docker
	@sudo docker run -p 8181:8080 -v /var/run/docker.sock:/var/run/docker.sock --rm $(IMAGE)
	#@PORT=1339 DEBUG=true ./$(EXECUTABLE)

docker-dev:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w $(LDFLAGS)'
	@sudo docker build --rm -t $(IMAGE) -f Dockerfile.dev .

dev: docker-dev
	@sudo docker run --rm -p 8181:8080 \
		-e "USERNAME=dacteev" \
		-e "PASSWORD=test" \
		-v $(shell pwd)/var/lib/docker-manager:/var/lib/docker-manager \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(shell pwd)/ui/dist:/opt/docker-manager/ui \
		--rm $(IMAGE)

stack-demo:
	#@sudo docker stack rm demo;
	@sudo docker stack deploy demo --compose-file demo/docker-compose.yml
