GLOBAL_BIN:=$(GOPATH)/bin

install-golangci-lint:
	GOBIN=$(GLOBAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.0

lint:
	GOBIN=$(GLOBAL_BIN) golangci-lint run ./... --config .golangci.pipeline.yaml