default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: build
build: ## go build
	go build .

.PHONY: install
install: build
	mkdir -p "/Users/harry.seeber/bin/plugins/registry.terraform.io/hashicorp/jamfpro"
	mv terraform-provider-jamfpro "/Users/harry.seeber/bin/plugins/registry.terraform.io/hashicorp/jamfpro/terraform-provider-jamfpro"
