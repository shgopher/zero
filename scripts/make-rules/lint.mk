# ==============================================================================
# Makefile helper functions for lintes
#

.PHONY: lint.run
lint.run: lint.ci lint.kubefiles lint.dockerfiles lint.charts ## Run all available linters.

.PHONY: lint.ci
lint.ci: lint.golngci-lint lint.zero ## Run CI-related linters.

.PHONY: lint.golangci-lint
lint.golangci-lint: tools.verify.golangci-lint ## Run golangci to lint source codes.
	@echo "===========> Run golangci to lint source codes"
	@golangci-lint run -c $(ZROOT)/.golangci.yaml $(ZROOT)/...

.PHONY: lint.zero
lint.zero: ## Run linters developed by zero developers.
	@$(GO) run cmd/lint-kubelistcheck/main.go $(ZROOT)/...

.PHONY: lint.kubefiles
lint.kubefiles: tools.verify.kube-linter ## Lint protobuf files.
	@kube-linter lint $(ZROOT)/deployments

.PHONY: lint.dockerfiles 
lint.dockerfiles: image.verify go.build.verify ## Lint dockerfiles.
	@$(ZROOT)/scripts/ci-lint-dockerfiles.sh $(HADOLINT_VER) $(HADOLINT_FAILURE_THRESHOLD)

.PHONY: lint.charts
lint.charts: tools.verify.helm ## Lint helm charts.
	$(MAKE) chart.lint
