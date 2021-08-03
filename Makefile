GIT_REVISION := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
VERSION := "$(shell git describe --tags --abbrev=0)-${GIT_REVISION}"
GIT_OPT := -X main.Branch=$(GIT_BRANCH) -X main.Revision=$(GIT_REVISION) -X main.Version=$(VERSION)
GO_OPT= -ldflags "$(GIT_OPT)"
RELEASE_OPT= -ldflags "-linkmode external -extldflags -static $(GIT_OPT)"

zombie:
	go build $(GO_OPT) -o ./bin/zombie main.go

dev-server:
	go build -o ./bin/dev-server ./cmd/dev/server.go

all: zombie dev-server

run:
	go run main.go --config zombie.yaml

test:
	go test ./... -v -count=1

release:
	go build $(RELEASE_OPT) -o ./bin/release/zombie -a main.go