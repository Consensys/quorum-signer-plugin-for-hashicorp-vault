OUTPUT_DIR := "$(shell pwd)/build"
NAME := "quorum-signer"
VERSION := "0.0.1-SNAPSHOT"

default: clean tools checkfmt test build

clean:
	@rm -rf ${OUTPUT_DIR}

checkfmt: tools
	@GO_FMT_FILES="$$(goimports -l `find . -name '*.go'`)"; \
	test -z "$${GO_FMT_FILES}" || ( echo "Please run 'make fixfmt' to format the following files: \n$${GO_FMT_FILES}"; exit 1 )

fixfmt: tools
	@goimports -w `find . -name '*.go'`

test:
	GOFLAGS="-mod=readonly" go test ./...

build:
	@mkdir -p ${OUTPUT_DIR}
	@echo Output to ${OUTPUT_DIR}
	@GOFLAGS="-mod=readonly" go build \
		-o "${OUTPUT_DIR}/${NAME}-${VERSION}" \
		.

tools: goimports

goimports:
ifeq (, $(shell which goimports))
	@GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
endif
