NAME=pg-chat-ops
BINARY=./bin/${NAME}
SOURCEDIR=./src
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

VERSION := $(shell git describe --abbrev=0 --tags)
SHA := $(shell git rev-parse --short HEAD)

GOPATH ?= /usr/local/go
GOPATH := ${CURDIR}:${GOPATH}
export GOPATH

$(BINARY): $(SOURCES)
	go build -o ${BINARY} -ldflags "-X main.BuildVersion=$(VERSION)-$(SHA)" $(SOURCEDIR)/$(NAME)/cmd/main.go

run: clean $(BINARY)
	${BINARY} --init-script ./examples/init.lua

test:
	go test -x -v ${NAME}/dsl

tar: clean
	mkdir -p rpm/SOURCES
	tar --transform='s,^\.,$(NAME)-$(VERSION),'\
		--exclude=rpm/SOURCES \
		-czf rpm/SOURCES/$(NAME)-$(VERSION).tar.gz .

docker: submodule_check tar
	/bin/cp -av $(CURDIR)/rpm /build
	/bin/cp -av $(CURDIR)/rpm/SPECS/rpm.spec /build/SPECS/$(NAME)-$(VERSION).spec
	sed -i 's|%define version unknown|%define version $(VERSION)|g' /build/SPECS/$(NAME)-$(VERSION).spec
	chown -R root:root /build
	rpmbuild -ba --define '_topdir /build'\
		/build/SPECS/$(NAME)-$(VERSION).spec

clean:
	rm -f $(BINARY)
	rm -f rpm-tmp.*

.DEFAULT_GOAL: $(BINARY)

include Makefile.git
