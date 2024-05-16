DEV      := deploymenttheory
PROVIDER := jamfpro
VERSION  := $(shell git describe --abbrev=0 --tags --match "v*")
PLUGINS  := ${HOME}/bin/plugins/registry.terraform.io/${DEV}/${PROVIDER}
BIN      := terraform-provider-jamfpro_${VERSION}

define TERRAFORMRC

add the following config to ~/.terraformrc to enable override:
```
provider_installation {
  dev_overrides {
    "${DEV}/${PROVIDER}" = "${PLUGINS}"
  }
}
```
endef

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Run go build. Output to dist/.
.PHONY: build
build:
	@mkdir -p dist
	go build -o dist/${BIN} .

# Run go build. Output to dist/.
.PHONY: build_override
build_override: build
	mkdir -p ${PLUGINS}
	mv dist/${BIN} ${PLUGINS}/${BIN}

# Run go build. Move artifact to terraform plugins dir. Output override config for ~/.terraformrc
.PHONY: install
install: build_override
	$(info ${TERRAFORMRC})
