.PHONY: swagger

# Directory where main application entry point is located
CMD_DIR := ./cmd

# Output directory for Swagger files
SWAGGER_OUTPUT_DIR := ./swagger

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	# Ensure the output directory exists
	@mkdir -p $(SWAGGER_OUTPUT_DIR)
	# Run the swag init command
	swag init -g $(CMD_DIR)/main.go --output $(SWAGGER_OUTPUT_DIR)
	@echo "Swagger documentation generated in $(SWAGGER_OUTPUT_DIR)/"
