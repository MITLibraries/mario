.PHONY: help install test tests update dist publish promote
SHELL=/bin/bash
ECR_REGISTRY=672626379771.dkr.ecr.us-east-1.amazonaws.com
DATETIME:=$(shell date -u +%Y%m%dT%H%M%SZ)

help: ## Print this message
	@awk 'BEGIN { FS = ":.*##"; print "Usage:  make <target>\n\nTargets:" } \
		/^[-_[:alpha:]]+:.?*##/ { printf "  %-15s%s\n", $$1, $$2 }' $(MAKEFILE_LIST)

generate: ## Generate code lists
	go generate

install: generate ## Install mario binary
	go install

test: generate ## Run tests
	go test -v ./...

tests: test

update: ## Update dependencies
	go get -u ./...

dist: ## Build docker image
	docker build -t $(ECR_REGISTRY)/mario-stage:latest \
		-t $(ECR_REGISTRY)/mario-stage:`git describe --always` \
		-t mario:latest .
	@tput setaf 2
	@tput bold
	@echo "Finished building docker image. Try running:"
	@echo "  $$ docker run --rm mario:latest"
	@tput sgr0

publish: dist ## Build, tag and push
	$$(aws ecr get-login --no-include-email --region us-east-1)
	docker push $(ECR_REGISTRY)/mario-stage:latest
	docker push $(ECR_REGISTRY)/mario-stage:`git describe --always`

promote: ## Promote the current staging build to production
	$$(aws ecr get-login --no-include-email --region us-east-1)
	docker pull $(ECR_REGISTRY)/mario-stage:latest
	docker tag $(ECR_REGISTRY)/mario-stage:latest $(ECR_REGISTRY)/mario-prod:latest
	docker tag $(ECR_REGISTRY)/mario-stage:latest $(ECR_REGISTRY)/mario-prod:$(DATETIME)
	docker push $(ECR_REGISTRY)/mario-prod:latest
	docker push $(ECR_REGISTRY)/mario-prod:$(DATETIME)
