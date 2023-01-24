.PHONY: *
default: build

# Run acceptance tests
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

build: fmt
	go install

fmt:
	go fmt

doc:
	go generate