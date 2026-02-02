.PHONY: swagger swagger2 swagger3 help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

swagger: swagger2 swagger3 ## Generate both Swagger 2.0 and OpenAPI 3.1 docs

swagger2: ## Generate Swagger 2.0 docs (for gin-swagger UI)
	@echo "Generating Swagger 2.0 docs..."
	@swag init -g ./cmd/api/main.go --output docs

swagger3: ## Generate OpenAPI 3.1 docs (for openapi-typescript)
	@echo "Generating OpenAPI 3.1 docs..."
	@swag init -g ./cmd/api/main.go --output docs/openapi3 --v3.1
