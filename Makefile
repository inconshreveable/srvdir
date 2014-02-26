.PHONY: default server client deps fmt clean all release-all release-server release-client contributors
export GOPATH:=$(shell pwd)

BUILDTAGS=debug
default: all

deps:
	go get -tags '$(BUILDTAGS)' -d -v srvdir/...

server: deps
	go install -gcflags "-N -l" -tags '$(BUILDTAGS)' srvdir/cmd/srvdird

fmt:
	go fmt srvdir/...

client: deps
	go install -gcflags "-N -l" -tags '$(BUILDTAGS)' srvdir/cmd/srvdir

release-client: BUILDTAGS=release
release-client: client

release-server: BUILDTAGS=release
release-server: server

release-all: fmt release-client release-server

all: fmt client server

clean:
	go clean -i -r srvdir/...

contributors:
	echo "Contributors to srvdir:\n" > CONTRIBUTORS
	git log --raw | grep "^Author: " | sort | uniq | cut -d ' ' -f2- | sed 's/^/- /' | cut -d '<' -f1 >> CONTRIBUTORS
