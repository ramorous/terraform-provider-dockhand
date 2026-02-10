.PHONY: build test fmt lint install doc testacc

build:
	go build -o terraform-provider-dockhand

test:
	go test -v ./...

fmt:
	gofmt -s -w .

lint:
	golangci-lint run

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/ramorous/dockhand/0.1.0/linux_amd64/
	cp terraform-provider-dockhand ~/.terraform.d/plugins/registry.terraform.io/ramorous/dockhand/0.1.0/linux_amd64/

doc:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate

clean:
	rm -f terraform-provider-dockhand
	go clean

testacc:
	TF_ACC=1 go test -v ./... -timeout 120m

.PHONY: all
all: build doc
