# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)


# ==============================================================================
# Hack

example1:
	go run api/cmd/example1/main.go

example2:
	go run api/cmd/example2/main.go

# ==============================================================================
# Install dependencies
#   https://ollama.com/
#   https://github.com/ollama/ollama/tree/main
#   https://github.com/tmc/langchaingo/

ollama:
	ollama run llama3

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor
