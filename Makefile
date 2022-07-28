OUTPUT_PATH=bin
BINARY=foto

VERSION=1.0.2
BUILD_TIME=`date +%Y%m%d%H%M`

LDFLAGS=-ldflags "-X github.com/waynezhang/foto/internal/cmd.Version=${VERSION} -X github.com/waynezhang/foto/internal/cmd.BuildTime=${BUILD_TIME}"

all: build

build:
	go build ${LDFLAGS} -o ${OUTPUT_PATH}/${BINARY} main.go

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${OUTPUT_PATH}/${BINARY} ] ; then rm ${OUTPUT_PATH}/${BINARY} ; fi
