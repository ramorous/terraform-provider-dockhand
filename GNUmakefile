default: install

.PHONY: build
build:
	go build -o terraform-provider-dockhand

.PHONY: release
release:
	@echo "Running GoReleaser..."
	goreleaser release --clean

.PHONY: release-snapshot
release-snapshot:
	@echo "Running GoReleaser in snapshot mode..."
	goreleaser release --snapshot --skip=publish --rm-dist

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/ramorous/dockhand/0.1.0/linux_amd64/
	cp terraform-provider-dockhand ~/.terraform.d/plugins/registry.terraform.io/ramorous/dockhand/0.1.0/linux_amd64/

.PHONY: uninstall
uninstall:
	rm -rf ~/.terraform.d/plugins/registry.terraform.io/ramorous/dockhand/

.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: testacc
testacc:
	TF_ACC=1 go test -v -race -cover ./...

.PHONY: fmt
fmt:
	gofmt -s -w .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: doc
doc:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate

.PHONY: generate
generate:
	go generate ./...

.PHONY: clean
clean:
	rm -f terraform-provider-dockhand
	rm -rf dist/
