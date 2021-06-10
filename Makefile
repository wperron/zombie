zombie:
	go build -o ./bin/zombie main.go

dev-server:
	go build -o ./bin/dev-server ./cmd/dev/server.go

all: zombie dev-server

run:
	go run main.go --config zombie.yaml

test:
	go test ./... -v -count=1