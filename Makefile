ver =

all: build

.PHONY: help server run build tag tidy release

help: ## Print the help menu
	@echo Usage: make [command]
	@echo
	@echo Commands:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf"  \033[36m%-30s\033[0m%s\n", $$1, $$2}'

server: ## Run the microservice locally
	go run cmd/dispatcher/main.go

run: build ## Run the microservice in a container
	docker run -p 6001:6001 -v $(shell pwd)/.env:/.env -it multimoml/dispatcher:latest

build: tidy ## Build the Docker image
	docker build -t multimoml/dispatcher:latest .

tag: ## Update the project version and create a Git tag with a changelog
    ifndef ver
		git tag -l
    else
		sed -i 's/:v[0-9.]*/:v'$(ver)'/' .github/workflows/publish.yml

		# Commit all changed files
		git add .

		# Creates a new Git tag with a changelog
		git commit -qm "Bump project version to $(ver)"
		printf "Release v$(ver)\n\nChangelog:\n" > changelog.txt
		git log $(shell git describe --tags --abbrev=0)..HEAD~1 --pretty=format:"  - %s" >> changelog.txt
		git tag -asF changelog.txt v$(ver)
		rm changelog.txt
    endif

tidy: ## Update dependencies
	go mod tidy

release: tag ## Create a new release and push it
	git fetch . main:prod
	git push --follow-tags origin main prod
