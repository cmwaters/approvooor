install:
	@echo "Installing BlobuSign"
	@go install ./cmd
.PHONY: install

build:
	@echo "Building BlobuSign"
	@go build -o build/blobusign ./cmd
.PHONY: build