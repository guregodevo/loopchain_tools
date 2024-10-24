# Plugin and version settings
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.2")
REPO := guregodevo/loopchain_tools # Update this to your repo
SRC_DIRS := yfinance_news yfinance python chrome example
OUTPUT_DIR := ./plugins

# Create the output directory if it doesn't exist
create_output_dir:
	@mkdir -p $(OUTPUT_DIR)

# Build Go plugins using Docker
build_plugins:
	@docker build -t go-plugin-builder:latest .
	@CONTAINER_ID=$$(docker create go-plugin-builder:latest) && \
	docker cp $$CONTAINER_ID:/app/plugins ./plugins && \
	docker rm $$CONTAINER_ID

# Create a new tag for release
tag_version:
	@if [ "$(VERSION)" = "" ]; then \
		echo "No tags found. Create a new tag before releasing."; \
		exit 1; \
	else \
		echo "Tagging version $(VERSION)"; \
		git tag -a $(VERSION) -m "Release $(VERSION)"; \
		git push origin $(VERSION); \
	fi

# Create a new release on GitHub and upload all the plugin files
release: build_plugins
	@echo "Releasing version $(VERSION)"
	@if gh release view $(VERSION) --repo $(REPO); then \
		echo "Release $(VERSION) already exists, uploading files..."; \
		for file in $(OUTPUT_DIR)/*.so; do \
			echo "Uploading $$file"; \
			gh release upload $(VERSION) $$file --repo $(REPO) --clobber; \
		done; \
	else \
		echo "Creating new release $(VERSION)"; \
		gh release create $(VERSION) $(OUTPUT_DIR)/*.so --repo $(REPO) --title "Release $(VERSION)" --notes "Automated release for version $(VERSION)" --target main; \
	fi
	@echo "Released all plugins for version $(VERSION) to GitHub."

# Clean up build artifacts
clean:
	@rm -rf $(OUTPUT_DIR)

# Ensure clean build before releasing
all: clean build_plugins release
