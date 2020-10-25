.PHONY: all build build-wasm clean test run

GO_BUILD_FLAGS=
GOROOT=/usr/local/Cellar/go/1.14.7/libexec


all: build-wasm

copy-support:
	echo ${GOROOT}
	cp "${GOROOT}/misc/wasm/wasm_exec.js" .

clean:
	go clean ./...

test:
	go test -race -p 1 -timeout 2m -v ./...

run:
	go run ${GO_BUILD_FLAGS} main.go

build: copy-support
	go build ${GO_BUILD_FLAGS} main.go

build-wasm: copy-support
	GOOS=js GOARCH=wasm go build ${GO_BUILD_FLAGS} -o main.wasm

