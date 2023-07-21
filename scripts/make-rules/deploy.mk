# ==============================================================================
# Makefile helper functions for deploy docker image in a test kubernetes
#

NAMESPACE ?= zero
CONTEXT ?= kind-zero

DEPLOYS=zero-usercenter zero-gateway

.PHONY: deploy.deploy
deploy.deploy: $(addprefix deploy.deploy., $(addprefix $(PLATFORM)., $(DEPLOYS))) ## Deploy all configured services.

.PHONY: deploy.deploy.%
deploy.deploy.%: image.build.% ## Deploy a specified service. (Note: Use `make deploy.<service>` to deploy a specific service.)
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM)))) 
	$(eval DEPLOY := $(word 2,$(subst ., ,$*)))
	@echo "===========> Deploying $(REGISTRY_PREFIX)/$(DEPLOY)-$(ARCH):$(VERSION)"
	@kubectl -n $(NAMESPACE) --context=$(CONTEXT) set image deployment/$(DEPLOY) $(DEPLOY)=$(REGISTRY_PREFIX)/$(DEPLOY)-$(ARCH):$(VERSION)

