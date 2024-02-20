install:
	@echo "Installing BlobuSign"
	@go install ./cmd
.PHONY: install

build:
	@echo "Building BlobuSign"
	@go build -o build/blobusign ./cmd
.PHONY: build

run-mock:
	@echo "Running BlobuSign Mock Server on :8080"
	@go run ./mock