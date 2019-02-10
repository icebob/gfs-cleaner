# SHELL=C:/Windows/System32/cmd.exe

# Get the current full sha from git
GITSHA:=$(shell git rev-parse HEAD)

# Get the current local branch name from git (if we can, this may be blank)
# GITBRANCH:=$(shell git symbolic-ref --short HEAD)

# BUILD_TIME:=`date +%FT%T%z`

# LDFLAGS for build
LDFLAGS=-ldflags="-X main.GitCommit=${GITSHA}"
# LDFLAGS=-ldflags="-X main.GitCommit=${GITSHA} -X main.BuildTime=${BUILD_TIME}"

# Build path style
BUILD_OUTPUT=-output="releases/{{.OS}}-{{.Arch}}/{{.Dir}}"


default: deps build

deps:
	go get -t ./...
	go get github.com/franciscocpg/gox

build: clean
	@gox -os="windows linux" -arch="386 amd64 arm" ${LDFLAGS} ${BUILD_OUTPUT} .

build-all: clean
	@gox ${LDFLAGS} ${BUILD_OUTPUT} .

clean:
	@rm -rf releases

install:
	@go install .

run:
	@go run ./main.go

test: vet
	echo "Missing tests"

packing:
	for f in releases/*/icepacker; do filename=$$(basename $$(dirname "$$f")); tar -cf "releases/icepacker-$$filename.tar.gz" -C $$(dirname $$f) $$(basename $$f) ; done; \
	for f in releases/*/icepacker.exe; do filename=$$(basename $$(dirname "$$f")); zip -j "releases/icepacker-$$filename.zip" $$f; done; \
