.PHONY: help install test tests update dist publish
SHELL=/bin/bash
ECR_REGISTRY=672626379771.dkr.ecr.us-east-1.amazonaws.com

help: ## Print this message
	@awk 'BEGIN { FS = ":.*##"; print "Usage:  make <target>\n\nTargets:" } \
		/^[-_[:alpha:]]+:.?*##/ { printf "  %-15s%s\n", $$1, $$2 }' $(MAKEFILE_LIST)

install: ## Install mario binary
	dep ensure
	go install

test: ## Run tests
	go test

tests: test

update: ## Update dependencies
	dep ensure -update

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
