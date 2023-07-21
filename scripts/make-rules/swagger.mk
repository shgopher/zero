# ==============================================================================
# Makefile helper functions for swagger
#

.PHONY: swagger.run
swagger.run: tools.verify.swagger
	@echo "===========> Generating swagger API docs"
	#@swagger generate spec --scan-models -w $(ZROOT)/cmd/gen-swagger-type-docs -o $(ZROOT)/api/swagger/kubernetes.yaml
	@swagger mixin `find $(ZROOT)/api/openapi -name "*.swagger.json"` \
		-q                                                    \
		--keep-spec-order                                     \
		--format=yaml                                         \
		--ignore-conflicts                                    \
		-o $(ZROOT)/api/swagger/swagger.yaml
	@echo "Generated at: $(ZROOT)/api/swagger/swagger.yaml"

.PHONY: swagger.serve
swagger.serve: tools.verify.swagger
	@swagger serve -F=redoc --no-open --port 65534 $(ZROOT)/api/swagger/swagger.yaml
