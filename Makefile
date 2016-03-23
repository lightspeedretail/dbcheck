# Important to know: By default, 'make' exits on tasks that return non-zero
# status codes.

all: clean bootstrap setup lint test ## Run ALL THE THINGS.

# In general, "make install" is used in projects where there's an actual
# artefact produced from the build installed to the executing system. Since
# we're pulling in dependencies for the build, "bootstrap" makes more sense.

bootstrap: ## Install dependencies.
	./script/bootstrap

test: bootstrap lint test-unit test-integration ## Runs lint rules, unit tests, and integration tests for the project.
    
test-unit: ## Runs the unit tests.
	./script/test-unit

lint: bootstrap ## Runs lint rules for Go.
	./script/lint

setup: bootstrap ## Set the project up to a working state.
	./script/setup

cibuild: lint test-unit ## Build target for the CI server.

clean: bootstrap ## Put the repository back in a pristine state.
	./script/clean

release: ## Tags a release and pushes the tag.
	./script/release

help: ## This file.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
.PHONY: bootstrap clean test lint setup cibuild release migrate test-integration test-unit all
