artifactName = libI18n

## which linter version to use?
GOLANGCI_LINT_VERSION = v1.24.0

buildDir = build
TOOLSDIR = build/tools
GO = go
GOLANGCI_LINT = ${TOOLSDIR}/golangci-lint
TMP_GOPATH = $(CURDIR)/${TOOLSDIR}/.gopath

all: generate lint test build ## Runs lint, test and build

clean: ## Removes any temporary and output files
	rm -rf ${buildDir}

lint: ## Executes all linters
	${GOLANGCI_LINT} run --enable-all

test: ## Executes the tests
	${GO} test -race ./...

.PHONY: build
build: ## Performs a build and puts everything into the build directory
	${GO} build -o ${buildDir}/${artifactName}

run: build ## Starts the compiled program
	${buildDir}/${artifactName}

help: ## Shows this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: generate
generate: ## Executes go generate
	${GO} generate

setup: installGolangCi ## Installs golint and golangci-lint
	${GO} mod tidy


installGolangCi:
	mkdir -p ${TOOLSDIR}
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLSDIR)/ $(GOLANGCI_LINT_VERSION)

.DEFAULT_GOAL := help

