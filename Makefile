.PHONY: dev build start clean

# Variables
BINARY_DIR=bin
CLI_BINARY=$(BINARY_DIR)/supercli


GOCMD=/usr/local/go/bin/go
GOBUILD=$(GOCMD) build

dev:
	@echo "💻 Running in development mode (hot reload with Air)..."
	air

build:
	@echo "🔨 Building CLI (supercli)..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(CLI_BINARY) ./main.go


clean:
	@echo "🧹 Cleaning..."
	rm -rf bin tmp

# Install CLI globally
install: build
	@echo "📦 Installing supercli..."
	cp $(CLI_BINARY) /usr/local/bin/supercli
	@echo "✅ minictl installed to /usr/local/bin/supercli"
#manual installation
# make build-cli
# sudo cp bin/supercli /usr/local/bin/supercli

# Uninstall CLI
uninstall:
	@echo "🗑️  Uninstalling supercli..."
	rm -f /usr/local/bin/supercli
	@echo "✅ supercli uninstalled"