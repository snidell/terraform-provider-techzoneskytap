TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
GOIMPORT_FILES?=$$(find . -type f -name '*.go' -not -path './vendor/*')
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=techzoneskytap

default: build

build: fmtcheck
	go build -o bin/terraform-provider-$(PKG_NAME)
	@sh -c "'$(CURDIR)/scripts/generate-dev-overrides.sh'"

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 240m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"


test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

lint:
	golint ./skytap/...

imports:
	goimports -w $(GOIMPORT_FILES)

BIN=$(CURDIR)/bin
$(BIN)/%:
	@echo Installing tools from tools.go
	@cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | GOBIN=$(BIN) xargs -tI {} go install {}

generate-docs: $(BIN)/tfplugindocs
	$(BIN)/tfplugindocs

tfproviderlint: $(BIN)/tfproviderlint
	$(BIN)/tfproviderlint $(TFPROVIDERLINT_ARGS) ./...

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test ./skytap -v -sweep=ALL $(SWEEPARGS) -timeout 30m


.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile lint imports generate-docs tfproviderlint sweep
