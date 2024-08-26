# Default target
.DEFAULT_GOAL := help

# Help
help:
	@echo "Get Started:"
	@echo "  download       - Download Go module dependencies"
	@echo "  tidy           - Tidy Go module dependencies"
	@echo ""
	@echo "Tag:"
	@echo "  tag            - Show current git tag"
	@echo "  tag-up         - Update git tag"
	@echo ""
	@echo "Helps:"
	@echo "  help           - Show this help message"
	@echo ""
	@echo "APP:"
	@echo "  build          - Build the project"
	@echo "  example            - Run the project as an example"
	@echo "  install            - Install the application"
	@echo "  test            	- Run tests"
	@echo ""
	@echo "Example:"
	@echo "  make example"
	@echo "  make install"



# Phony targets
.PHONY: help
