PROJ_NAME = gnmatcher

VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S%Z')

FLAGS_SHARED =  $(NO_C) GOARCH=amd64
FLAGS_LINUX = $(FLAGS_SHARED) GOOS=linux
FLAGS_MAC = $(FLAGS_SHARED) GOOS=darwin
FLAGS_WIN = $(FLAGS_SHARED) GOOS=windows
NO_C = CGO_ENABLED=0

FLAGS_LD = -ldflags "-X github.com/gnames/$(PROJ_NAME)/pkg.Build=$(DATE) \
                  -X github.com/gnames/$(PROJ_NAME)/pkg.Version=$(VERSION)"
FLAGS_REL = -trimpath -ldflags "-s -w \
						-X github.com/gnames/$(PROJ_NAME)/pkg.Build=$(DATE)"

GOCMD= go
GOINSTALL = $(GOCMD) install $(FLAGS_LD)
GOBUILD = $(GOCMD) build $(FLAGS_LD)
GORELEASE = $(GOCMD) build $(FLAGS_REL)
GOCLEAN  = $(GOCMD) clean
GOGET = $(GOCMD) get

all: install

test: deps install
	go test -parallel=1 -shuffle=on -count=1 -race -coverprofile=coverage.txt -covermode=atomic ./...

tools: deps
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

deps:
	@echo Download go.mod dependencies
	$(GOCMD) mod download;

build:
	$(GOCLEAN); \
	$(NO_C) $(GOBUILD);

dc: build
	docker-compose build;

buildrel:
	$(GOCLEAN); \
	$(NO_C) $(GORELEASE);

release: dockerhub
	$(GOCLEAN); \
	$(FLAGS_LINUX) $(GORELEASE); \
	tar zcvf /tmp/$(PROJ_NAME)-$(VER)-linux-x86.tar.gz $(PROJ_NAME); \
	$(GOCLEAN);

install:
	$(NO_C) $(GOINSTALL);

docker:
	$(GOCLEAN); \
	$(FLAGS_LINUX) $(GORELEASE); \
	docker buildx build -t gnames/$(PROJ_NAME):latest -t gnames/$(PROJ_NAME):$(VERSION) .; \

dockerhub: docker
	docker push gnames/$(PROJ_NAME); \
	docker push gnames/$(PROJ_NAME):$(VERSION)

