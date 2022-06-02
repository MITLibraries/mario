.PHONY: help install test tests update dist publish promote
SHELL=/bin/bash
DATETIME:=$(shell date -u +%Y%m%dT%H%M%SZ)
### This is the Terraform-generated header for mario-dev
ECR_NAME_DEV:=mario-dev
ECR_URL_DEV:=222053980223.dkr.ecr.us-east-1.amazonaws.com/mario-dev
### End of Terraform-generated header ###

help: ## Print this message
	@awk 'BEGIN { FS = ":.*##"; print "Usage:  make <target>\n\nTargets:" } \
		/^[-_[:alpha:]]+:.?*##/ { printf "  %-15s%s\n", $$1, $$2 }' $(MAKEFILE_LIST)

install: ## Install mario binary
	go install ./...

test: ## Run tests
	go test -v ./...

tests: test

update: ## Update dependencies
	go get -u ./...

### Developer Deploy Commands ###
dist-dev: ## Build docker container (intended for developer-based manual build)
	docker build --platform linux/amd64 \
	    -t $(ECR_URL_DEV):latest \
		-t $(ECR_URL_DEV):`git describe --always` \
		-t $(ECR_NAME_DEV):latest .

publish-dev: dist-dev ## Build, tag and push (intended for developer-based manual publish)
	docker login -u AWS -p $$(aws ecr get-login-password --region us-east-1) $(ECR_URL_DEV)
	docker push $(ECR_URL_DEV):latest
	docker push $(ECR_URL_DEV):`git describe --always`