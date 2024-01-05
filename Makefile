OUTPUT_PATH=bin
BINARY=foto

VERSION=`git for-each-ref --sort=creatordate --format '%(refname)' refs/tags | tail -n 1 | sed 's/refs\/tags\/v\(.*\)/\1/g'`
BUILD_TIME=`date +%Y%m%d%H%M`

LDFLAGS=-ldflags "-X github.com/waynezhang/foto/internal/cmd.Version=${VERSION} -X github.com/waynezhang/foto/internal/cmd.Revision=${BUILD_TIME}"

all: build

build:
	@go build ${LDFLAGS} -o ${OUTPUT_PATH}/${BINARY} main.go

test:
	@go test ./...

coverage:
	@TMPFILE=$$(mktemp); \
		go test ./... -coverprofile=$$TMPFILE; \
		go tool cover -html $$TMPFILE

changelog:
	@whatchanged HEAD > CHANGELOG.md

.PHONY: install
install:
	@go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	@if [ -f ${OUTPUT_PATH}/${BINARY} ] ; then rm ${OUTPUT_PATH}/${BINARY} ; fi
