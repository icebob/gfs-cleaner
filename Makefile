NAME=gfs-cleaner
VERSION=$(shell cat VERSION)
GITSHA=$(shell git rev-parse --short HEAD)

EXT_LD_FLAGS="-Wl,--allow-multiple-definition"
LD_FLAGS="-w -X main.Version=$(VERSION) -X main.GitCommit=${GITSHA} -extldflags=$(EXT_LD_FLAGS)"

default: deps build

deps:
	go get -t ./...

build:
	go mod download
	CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o ${NAME}

build-all: clean
	mkdir -p releases
	GOOS=darwin  GOARCH=amd64 go build -tags release -ldflags $(LD_FLAGS) -o releases/${NAME}-$(VERSION)-darwin-amd64
	GOOS=linux   GOARCH=386   go build -tags release -ldflags $(LD_FLAGS) -o releases/${NAME}-$(VERSION)-linux-386
	GOOS=linux   GOARCH=amd64 go build -tags release -ldflags $(LD_FLAGS) -o releases/${NAME}-$(VERSION)-linux-amd64
	GOOS=linux   GOARCH=arm   go build -tags release -ldflags $(LD_FLAGS) -o releases/${NAME}-$(VERSION)-linux-arm
	GOOS=linux   GOARCH=arm64 go build -tags release -ldflags $(LD_FLAGS) -o releases/${NAME}-$(VERSION)-linux-arm64
	GOOS=windows GOARCH=386   go build -tags release -ldflags $(LD_FLAGS) -o releases/${NAME}-$(VERSION)-windows-386.exe
	GOOS=windows GOARCH=amd64 go build -tags release -ldflags $(LD_FLAGS) -o releases/${NAME}-$(VERSION)-windows-amd64.exe
	cd releases; sha256sum * > sha256sums.txt

clean:
	@rm -rf releases

install:
	@go install .

run:
	@go run ./main.go clean test --dry

test:
	echo "Missing tests"
