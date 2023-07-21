# ==============================================================================
#  Makefile helper functions for tools
#

# Specify tools category.
CODE_GENERATOR_TOOLS= client-gen lister-gen informer-gen defaulter-gen deepcopy-gen prerelease-lifecycle-gen conversion-gen openapi-gen
# code-generator is a makefile target not a real tool.
CI_WORKFLOW_TOOLS := code-generator golangci-lint goimports wire 
# unused tools in this project: gentool
OTHER_TOOLS := mockgen gsemver git-chglog addlicense kratos kind go-apidiff gotests cfssl go-gitlint kustomize kafkactl kube-linter \
							 kubeconform kubectl helm-docs db2struct gentool air swagger license gothanks kubebuilder go-junit-report controller-gen
MANUAL_INSTALL_TOOLS := helm kafka

.PHONY: tools.install
tools.install: install.ci _install.other tools.print-manual-tool ## Install all tools.

.PHONY: install.ci
install.ci: $(addprefix tools.install., $(CI_WORKFLOW_TOOLS)) ## Install necessary tools used by CI/CD workflow.

.PHONY: _install.other
_install.other: $(addprefix tools.install., $(OTHER_TOOLS))

.PHONY: tools.print-manual-tool
tools.print-manual-tool: 
	@echo "===========> The following tools may need to be installed manually:"
	@echo $(MANUAL_INSTALL_TOOLS) | awk 'BEGIN{RS=" "} {printf("%15s%s\n","- ",$$0)}'

.PHONY: tools.install.%
tools.install.%: ## Install a specified tool.
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

.PHONY: tools.verify.%
tools.verify.%: ## Verify a specified tool.
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

.PHONY: tools.verify.code-generator
tools.verify.code-generator: $(addprefix tools.verify., $(CODE_GENERATOR_TOOLS)) ## Verify a specified tool.

.PHONY: install.code-generator
install.code-generator: $(addprefix tools.install.code-generator., $(CODE_GENERATOR_TOOLS)) ## Install all necessary code-generator tools.

.PHONY: install.code-generator.%
install.code-generator.%: ## Install specified code-generator tool.
	#@$(GO) install k8s.io/code-generator/cmd/$*@$(CODE_GENERATOR_VERSION)
	@$(GO) install github.com/colin404/code-generator/cmd/$*@$(CODE_GENERATOR_VERSION)

.PHONY: install.swagger
install.swagger:
	@$(GO) install github.com/go-swagger/go-swagger/cmd/swagger@latest

.PHONY: install.golangci-lint
install.golangci-lint: ## Install golangci-lint.
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
	@$(ZROOT)/scripts/add-completion.sh golangci-lint bash

.PHONY: install.go-junit-report
install.go-junit-report:
	@$(GO) install github.com/jstemmer/go-junit-report@latest

.PHONY: install.wire
install.wire: ## Install wire.
	@$(GO) install github.com/google/wire/cmd/wire@latest

.PHONY: install.mockgen
install.mockgen: ## Install mockgen.
	@$(GO) install github.com/golang/mock/mockgen@latest

.PHONY: install.gotests
install.gotests: ## Install gotests.
	@$(GO) install github.com/cweill/gotests/gotests@latest

.PHONY: install.goimports
install.goimports: ## Install goimports.
	@$(GO) install golang.org/x/tools/cmd/goimports@latest

.PHONY: install.go-gitlint
install.go-gitlint: ## Install go-gitlint.
	@$(GO) install github.com/marmotedu/go-gitlint/cmd/go-gitlint@latest

.PHONY: install.gsemver
install.gsemver: ## Install gsemver.
	@$(GO) install github.com/arnaud-deprez/gsemver@latest
	@$(ZROOT)/scripts/add-completion.sh gsemver bash

.PHONY: install.git-chglog
install.git-chglog: ## Install git-chglog tool which is used to generate CHANGELOG.
	@$(GO) install github.com/git-chglog/git-chglog/cmd/git-chglog@latest

.PHONY: install.cfssl
install.cfssl: ## Install cfssl toolkit.
	@$(ZROOT)/scripts/install.sh zero::install::install_cfssl

.PHONY: install.addlicense
install.addlicense: ## Install addlicense.
	@$(GO) install github.com/superproj/addlicense@latest

.PHONY: install.kustomize
install.kustomize: ## Install kustomize.
	@$(GO) install sigs.k8s.io/kustomize/kustomize/v5@latest
	@$(ZROOT)/scripts/add-completion.sh kustomize bash

.PHONY: install.controller-gen
install.controller-gen: ## Install controller-gen.
	@$(GO) install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: install.kind
install.kind: ## Install kind cluster command line tool.
	@$(GO) install sigs.k8s.io/kind@latest
	@$(ZROOT)/scripts/add-completion.sh kind bash

.PHONY: install.go-apidiff
install.go-apidiff: ## Install go-apidiff.
	@$(GO) install github.com/joelanford/go-apidiff@latest

.PHONY: install.helm
install.helm: ## Install helm command line tool.
	@echo "does not implement yet"
	@exit 1
	@$(ZROOT)/scripts/add-completion.sh helm bash

.PHONY: install.kratos
install.kratos: ## Install kratos toolkit, includes multiple protoc plugins.
	@$(GO) install github.com/joelanford/go-apidiff@latest
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@$(GO) install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	@$(GO) install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	@$(GO) install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	@$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@$(GO) install github.com/envoyproxy/protoc-gen-validate@latest
	@$(GO) install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	@$(ZROOT)/scripts/add-completion.sh kratos bash

.PHONY: install.buf
install.buf: ## Install buf command line tool.
	@$(GO) install github.com/bufbuild/buf/cmd/buf@latest

.PHONY: install.kafkactl
install.kafkactl: ## Install kafkactl command line tool.
	@$(GO) install github.com/deviceinsight/kafkactl@latest
	@$(ZROOT)/scripts/add-completion.sh kafkactl bash

# kube-linter reference: https://docs.kubelinter.io/#/
.PHONY: install.kube-linter
install.kube-linter: ## Install kube-linter command line tool.
	@$(GO) install golang.stackrox.io/kube-linter/cmd/kube-linter@latest
	@$(ZROOT)/scripts/add-completion.sh kube-linter bash

.PHONY: install.kubeconform
install.kubeconform: ## Install kubeconform command line tool.
	@$(GO) install github.com/yannh/kubeconform/cmd/kubeconform@latest

.PHONY: install.kubectl
install.kubectl: ## Install kubectl command line tool.
	@curl --create-dirs -L -o $$HOME/bin/kubectl "https://dl.k8s.io/release/$(shell curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
	@chmod +x $$HOME/bin/kubectl
	@$(ZROOT)/scripts/add-completion.sh kubectl bash

.PHONY: install.helm-docs
install.helm-docs: ## Install helm-docs which is a tool to generating markdown documentation for helm charts.
	@$(GO) install github.com/norwoodj/helm-docs/cmd/helm-docs@latest

.PHONY: install.gentool
install.gentool: ## Install gentool which is a tool used to generate gorm model and query code.
	@$(GO) install gorm.io/gen/tools/gentool@latest

# db2struct --gorm --json -H 127.0.0.1 -d zero -t secret --package model --struct SecretM -u gateway -p 'proj(#)666' --target=secret.go
.PHONY: install.db2struct
install.db2struct: ## Install db2struct which is a tool used to converts a mysql table into a golang struct.
	@$(GO) install github.com/Shelnutt2/db2struct/cmd/db2struct@latest

.PHONY: install.protoc-go-inject-tag
install.protoc-go-inject-tag:
	@$(GO) install github.com/favadi/protoc-go-inject-tag@latest

.PHONY: install.air
install.air: ## Install air tool which is used to live reload your go apps.
	@$(GO) install github.com/cosmtrek/air@latest

.PHONY: install.license 
install.license : ## Install license tool which is used to generate LICENSE file as you want.
	@$(GO) install github.com/nishanths/license/v5@latest

.PHONY: install.gothanks
install.gothanks: ## Install gothanks tool which is used to automatically stars your go.mod github dependencies.
	@$(GO) install github.com/psampaz/gothanks@latest

.PHONY: install.kubebuilder
install.kubebuilder : ## Install kubebuilder tool which is used to building Kubernetes APIs using CRDs.
	# download kubebuilder and install locally.
	@curl -sL -o kubebuilder https://go.kubebuilder.io/dl/latest/$(shell $(GO) env GOOS)/$(shell $(GO) env GOARCH)
	@mkdir -p ${HOME}/bin
	@chmod +x kubebuilder && mv kubebuilder ${HOME}/bin

# gomodifytags -all -add-tags json -w -transform camelcase --skip-unexported -file *.go
.PHONY: install.gomodifytags
install.gomodifytags: ## Install gomodifytags tool which is used to modify struct field tags.
	@$(GO) install github.com/fatih/gomodifytags@latest
