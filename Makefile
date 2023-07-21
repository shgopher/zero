# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
#all: format tidy gen add-copyright lint cover build
all: format tidy gen build

# ==============================================================================
# Includes

include scripts/make-rules/common.mk # make sure include common.mk at the first include line
include scripts/make-rules/all.mk

# ==============================================================================
# Usage

define USAGE_OPTIONS

\033[35mOptions:\033[0m
  DEBUG            Whether to generate debug symbols. Default is 0.
  BINS             The binaries to build. Default is all of cmd.
                   This option is available when using: make build/build.multiarch
                   Example: make build BINS="zero-apiserver zero-miner-controller"
  IMAGES           Backend images to make. Default is all of cmd starting with zero-.
                   This option is available when using: make image/image.multiarch/push/push.multiarch
                   Example: make image.multiarch IMAGES="zero-apiserver zero-miner-controller"
  DEPLOYS          Deploy all configured services.
  REGISTRY_PREFIX  Docker registry prefix. Default is superproj. 
                   Example: make push REGISTRY_PREFIX=ccr.ccs.tencentyun.com/superproj VERSION=v1.6.2
  PLATFORMS        The multiple platforms to build. Default is linux_amd64 and linux_arm64.
                   This option is available when using: make build.multiarch/image.multiarch/push.multiarch
                   Example: make image.multiarch IMAGES="zero-apiserver zero-miner-controller" PLATFORMS="linux_amd64 linux_arm64"
  PLATFORMS        The multiple platforms to build. Default is linux_amd64 and linux_arm64.
                   This option is available when using: make build.multiarch/image.multiarch/push.multiarch
  MULTISTAGE       Set to 1 to build docker images using multi-stage builds. Default is 0.
  VERSION          The version information compiled into binaries.
                   The default is obtained from gsemver or git.
  A                Run all similar targets. Default only run CI-related rules.
  V                Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

## --------------------------------------
## Generate / Manifests
## --------------------------------------

##@ Generate

.PHONY: gen
gen: generated-files ## Generate CI-related files. Generate all files by specifying `A=1`.
	$(MAKE) gen.run
	if [[ "$(A)" == 1 ]]; then                                             \
		$(MAKE) gen.docgo.doc gen.appdocs gen.ca gen.kubeconfig ;                        \
	fi

.PHONY: generated-files
generated-files: ## Generate all necessary kubernetes related files, such as deepcopy files
	$(MAKE) -s generated.files

.PHONY: protoc
protoc: ## Generate api proto files.
	$(MAKE) gen.protoc

.PHONY: ca
ca: ## Generate CA files for all zero components.
	$(MAKE) gen.ca

## --------------------------------------
## Binaries
## --------------------------------------

##@ Build

.PHONY: build
build: tidy ## Build source code for host platform.
	$(MAKE) go.build

.PHONY: build.multiarch
build.multiarch: ## Build source code for multiple platforms. See option PLATFORMS.
	$(MAKE) go.build.multiarch

.PHONY: image
image: ## Build docker images for host arch.
	$(MAKE) image.build

.PHONY: image.multiarch
image.multiarch: ## Build docker images for multiple platforms. See option PLATFORMS.
	$(MAKE) image.build.multiarch

.PHONY: push
push: ## Build docker images for host arch and push images to registry.
	$(MAKE) image.push

.PHONY: push.multiarch
push.multiarch: ## Build docker images for multiple platforms and push images to registry.
	$(MAKE) image.push.multiarch

.PHONY: deploy
deploy: ## Build docker images for host arch.
	$(MAKE) deploy.deploy

## --------------------------------------
## Cleanup
## --------------------------------------

##@ Clean

.PHONY: clean
clean: ## Remove all files that are created by building and generaters.
	@echo "===========> Cleaning all build output and generated files"
	@-rm -vrf $(OUTPUT_DIR)
	@-rm -vrf $(ZROOT)/pkg/generated
	@-rm -vrf $(META_DIR)
	@-rm -vrf $(GENERATED_DOCKERFILE_DIR)
	@find $(APIROOT) -type f -regextype posix-extended -regex ".*.swagger.json|.*.pb.go" -delete
	@$(FIND) -type f -name 'zz_generated.*go' -delete
	@$(FIND) -type f -name '*_generated.go' -delete
	@$(FIND) -type f -name 'types_swagger_doc_generated.go' -delete

## --------------------------------------
## Testing
## --------------------------------------

##@ Test

.PHONY: test
test: ## Run unit test.
	$(MAKE) go.test

.PHONY: cover 
cover: ## Run unit test and get test coverage.
	$(MAKE) go.test.cover

## --------------------------------------
## Lint / Verification
## --------------------------------------

##@ Lint and Verify

.PHONY: lint
lint: ## Run CI-related linters. Run all linters by specifying `A=1`.
	if [[ "$(A)" == 1 ]]; then                                             \
		$(MAKE) lint.run ;                                                 \
	else ;                                                                \
	    $(MAKE) lint.ci ;                                                  \
	fi

.PHONY: apidiff
apidiff: tools.verify.go-apidiff ## Run the go-apidiff to verify any API differences compared with origin/master.
	@go-apidiff master --compare-imports --print-compatible --repo-path=.

## --------------------------------------
## Hack / Tools
## --------------------------------------

##@ Hack and Tools

.PHONY: format
format: tools.verify.goimports ## Run CI-related formaters. Run all formaters by specifying `A=1`.
	@echo "===========> Formating codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(PRJ_SRC_PATH)
	@$(GO) mod edit -fmt
	if [[ "$(A)" == 1 ]]; then                                             \
		$(MAKE) format.protobuf ;                                           \
	fi

.PHONY: format.protobuf
format.protobuf: tools.verify.buf ## Lint protobuf files.
	@for f in $(shell find $(APIROOT) -name *.proto) ; do                  \
	  buf format -w $$f ;                                                  \
	done

.PHONY: add-copyright
add-copyright: ## Ensures source code files have copyright license headers.
	$(MAKE) copyright.add

.PHONY: swagger
#swagger: gen.protoc
swagger: ## Generate and aggregate swagger document.
	@$(MAKE) swagger.run

.PHONY: swagger.serve
serve-swagger: ## Serve swagger spec and docs at 65534.
	@$(MAKE) swagger.serve

.PHONY: tidy
tidy:
	@$(GO) mod tidy

.PHONY: air.%
air.%: tools.verify.air
	@air -build.cmd='make build BINS=zero-$*' -build.bin='$(OUTPUT_DIR)/platforms/$(shell go env GOOS)/$(shell go env GOARCH)/zero-$*'

.PHONY: install-tools
install-tools: ## Install CI-related tools. Install all tools by specifying `A=1`.
	$(MAKE) install.ci
	if [[ "$(A)" == 1 ]]; then                                             \
		$(MAKE) _install.other ;                                            \
	fi

.PHONY: targets
targets: Makefile ## Show all Sub-makefile targets.
	@for mk in `echo $(MAKEFILE_LIST) | sed 's/Makefile //g'`; do echo -e \\n\\033[35m$$mk\\033[0m; awk -F':.*##' '/^[0-9A-Za-z._-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^\$$\([0-9A-Za-z_-]+\):.*?##/ { gsub("_","-", $$1); printf "  \033[36m%-45s\033[0m %s\n", tolower(substr($$1, 3, length($$1)-7)), $$2 }' $$mk;done;

.PHONY: help
help: Makefile ## Show this help info.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<TARGETS> <OPTIONS>\033[0m\n\n\033[35mTargets:\033[0m\n"} /^[0-9A-Za-z._-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^\$$\([0-9A-Za-z_-]+\):.*?##/ { gsub("_","-", $$1); printf "  \033[36m%-45s\033[0m %s\n", tolower(substr($$1, 3, length($$1)-7)), $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' Makefile #$(MAKEFILE_LIST)
	@echo -e "$$USAGE_OPTIONS"
