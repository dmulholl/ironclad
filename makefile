.PHONY: help
help: ## Prints available make commands.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / \
	{printf "\033[1;36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: all
all: ## Builds all command binaries.
	@rm -rf ./build/** && mkdir -p ./build
	@set -e && \
		for cmd in $$(find cmd/* -maxdepth 0 -type d | xargs -I {} basename {}); do \
			printf "\e[1;32mBuilding\e[0m cmd/$$cmd\n" && \
			go build -o ./build/$$cmd ./cmd/$$cmd; \
		done

.PHONY: build
build: ## Builds the primary application binary.
	go build -o ./build/ironclad ./cmd/ironclad

.PHONY: install
install: ## Builds and installs the primary application binary.
	go install ./cmd/ironclad

.PHONY: test
test: ## Runs unit tests.
	go test ./...

.PHONY: test-verbose
test-verbose: ## Runs unit tests verbosely.
	go test ./... -v

.PHONY: clean
clean: ## Deletes all build artifacts.
	@rm -rf ./build

.PHONY: release
release: ## Builds a set of zipped release binaries.
release: clean
	@printf "\e[1;32mCompiling\e[0m ...\n"
	GOOS=linux GOARCH=amd64 go build -o build/ironclad-linux-amd64/ironclad ./cmd/ironclad
	GOOS=darwin GOARCH=amd64 go build -o build/ironclad-mac-amd64/ironclad ./cmd/ironclad
	GOOS=darwin GOARCH=arm64 go build -o build/ironclad-mac-arm64/ironclad ./cmd/ironclad
	GOOS=windows GOARCH=amd64 go build -o build/ironclad-windows-amd64/ironclad.exe ./cmd/ironclad
	@printf "\n\e[1;32mBinaries\e[0m ...\n"
	@tree build
	@printf "\n\e[1;32mSHA 256 checksums\e[0m ...\n"
	@shasum -a 256 build/*/*
	@printf "\n\e[1;32mZipping\e[0m ...\n"
	@mkdir -p build/zipped
	@cd build && zip -r ./zipped/ironclad-linux-amd64.zip ironclad-linux-amd64 > /dev/null
	@cd build && zip -r ./zipped/ironclad-mac-amd64.zip ironclad-mac-amd64 > /dev/null
	@cd build && zip -r ./zipped/ironclad-mac-arm64.zip ironclad-mac-arm64 > /dev/null
	@cd build && zip -r ./zipped/ironclad-windows-amd64.zip ironclad-windows-amd64 > /dev/null
	@tree build/zipped
