OUTPUT_DIR := "$(shell pwd)/build"
NAME := "quorum-signer"
VERSION := "0.2.1-SNAPSHOT"
OS_ARCH := "$(shell go env GOOS)-$(shell go env GOARCH)"
BUILD_LD_FLAGS=-s -w $(extraldflags)

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
		-ldflags='$(BUILD_LD_FLAGS)' \
		-o "${OUTPUT_DIR}/dist/${NAME}-${VERSION}-${OS_ARCH}" \
		.

package: build
	@shasum -a 256 ${OUTPUT_DIR}/dist/${NAME}-${VERSION}-${OS_ARCH} | awk '{print $$1}' > ${OUTPUT_DIR}/dist/${NAME}-${VERSION}-${OS_ARCH}.checksum
	@zip -j -FS -q ${OUTPUT_DIR}/dist/${NAME}-${VERSION}-${OS_ARCH}.zip ${OUTPUT_DIR}/dist/*
	@shasum -a 256 ${OUTPUT_DIR}/dist/${NAME}-${VERSION}-${OS_ARCH}.zip | awk '{print $$1}' > ${OUTPUT_DIR}/dist/${NAME}-${VERSION}-${OS_ARCH}.zip.checksum

tools: goimports

goimports:
ifeq (, $(shell which goimports))
	@GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
endif
