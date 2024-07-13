OUTPUT_DIR=./bin
GITHUB_ACCOUNT=andyfusniak
GIT_COMMIT=$(shell git rev-parse --short HEAD)
VERSION=v0.1.0-dev

all: sitebuild

sitebuild:
	@go build -o $(OUTPUT_DIR)/sitebuild -ldflags "-X 'main.version=${VERSION}' -X 'main.gitCommit=${GIT_COMMIT}'" ./main.go

sitebuild-darwin-amd64:
	@GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/sitebuild -ldflags "-X 'main.version=${VERSION}' -X 'main.gitCommit=${GIT_COMMIT}'" ./main.go

.PHONY: clean
clean:
	-@rm -r $(OUTPUT_DIR)/sitebuild 2> /dev/null || true
