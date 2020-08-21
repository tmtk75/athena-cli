VERSION := $(shell git describe --tags --abbrev=0)
COMMIT:= $(shell git rev-parse HEAD)
VAR_VERSION := main.Version
VAR_COMMIT:= main.Commit

LDFLAGS := -ldflags "-X $(VAR_VERSION)=$(VERSION) \
	-X $(VAR_COMMIT)=$(COMMIT)"

athena-cli: *.go
	go build $(LDFLAGS) -o athena-cli .

release: build
	ghr -u tmtk75 $(VERSION) ./build

build: build/sha256sum.txt

build/sha256sum.txt: build/athena-cli_darwin_amd64.zip build/athena-cli_linux_amd64.zip
	(cd build && sha256sum athena-cli_* > sha256sum.txt)

build/athena-cli_darwin_amd64.zip build/athena-cli_linux_amd64.zip: build/athena-cli_darwin_amd64 build/athena-cli_linux_amd64
	parallel '(cd build && zip -m athena-cli_{1}_amd64.zip athena-cli_{1}_amd64)' \
	  ::: darwin linux

build/athena-cli_darwin_amd64 build/athena-cli_linux_amd64: *.go
	parallel 'GOARCH=amd64 GOOS={1} go build $(LDFLAGS) -o build/athena-cli_{1}_amd64 .' \
	  ::: darwin linux

clean:
	rm -rf athena-cli

distclean: clean
	rm -rf build
