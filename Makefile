VERSION := $(shell git describe --tags --abbrev=0)
COMMIT:= $(shell git rev-parse HEAD)
VAR_VERSION := main.Version
VAR_COMMIT:= main.Commit

LDFLAGS := -ldflags "-X $(VAR_VERSION)=$(VERSION) \
	-X $(VAR_COMMIT)=$(COMMIT)"

build: athena-cli

athena-cli: *.go
	go build $(LDFLAGS) -o athena-cli .

install:
	go install $(LDFLAGS)

clean:
	rm -rf athena-cli
