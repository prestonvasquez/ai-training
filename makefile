# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)


# ==============================================================================
# Hack

example1:
	go run examples/example1/main.go

example2:
	go run examples/example2/main.go

example3:
	go run -exec "env DYLD_LIBRARY_PATH=$$GOPATH/src/github.com/ardanlabs/vector/foundation/word2vec/libw2v/lib" examples/example3/main.go

example4:
	go run examples/example4/main.go

# ==============================================================================
# Install dependencies
#   https://ollama.com/
#   https://github.com/ollama/ollama/tree/main
#   https://github.com/tmc/langchaingo/

docker:
	docker pull mongodb/mongodb-atlas-local
	docker pull ollama/ollama

dev-up:
	docker-compose -f zarf/docker/compose.yaml up

dev-down:
	docker-compose -f zarf/docker/compose.yaml down

download-data:
	curl -o zarf/data/example3.gz -X GET http://snap.stanford.edu/data/amazon/productGraph/categoryFiles/reviews_Cell_Phones_and_Accessories_5.json.gz \
	&& gunzip -k -d zarf/data/example3.gz \
	&& mv zarf/data/example3 zarf/data/example3.json

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor
