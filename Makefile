# Plugin and version settings
PLUGIN_NAME := example_tool
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.1")
PLUGIN_SO := $(PLUGIN_NAME).so
REPO := guregodevo/loopchain_tools # Update this to your repo

# Build the plugin
build_plugin:
	go build -buildmode=plugin -o $(PLUGIN_SO)

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

# Create a new release on GitHub and upload the plugin file
release: build_plugin
	@gh release create $(VERSION) $(PLUGIN_SO) --repo $(REPO) --title "Release $(VERSION)" --notes "Automated release for version $(VERSION)" --target main
	@echo "Released $(VERSION) to GitHub."

# Clean up build artifacts
clean:
	rm -f $(PLUGIN_SO)

# Ensure clean build before releasing
all: clean build_plugin release
