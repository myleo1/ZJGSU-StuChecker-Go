BINARY=ZJGSU-StuChecker-Go
VERSION=v1.0.0-beta
DATE=`date +%FT%T%z`
GoVersion=`go version`
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION} -X 'main.date=${DATE}' -X 'main.goVersion=${GoVersion}'"
.PHONY: build build_osx fmt

default:
	@echo ${BINARY}
	@echo ${VERSION}
	@echo ${DATE}
	@echo ${GoVersion}

build:
	@GOOS=linux GOARCH=amd64 go build -trimpath -o build/${BINARY} ${LDFLAGS}
	@echo "[ok] build ${BINARY}"

build_osx:
	@go build -trimpath -o build/${BINARY} ${LDFLAGS}
	@echo "[ok] build_osx"

fmt:
	@gofmt -s -w ./