.PHONY: help install test tests update dist publish promote
SHELL=/bin/bash
ECR_REGISTRY=222053980223.dkr.ecr.us-east-1.amazonaws.com
# old org prod 672626379771.dkr.ecr.us-east-1.amazonaws.com
# new org dev 222053980223.dkr.ecr.us-east-1.amazonaws.com
DATETIME:=$(shell date -u +%Y%m%dT%H%M%SZ)

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

dist: ## Build docker image
	docker build -t $(ECR_REGISTRY)/dev-timdex-mario:latest \
		-t $(ECR_REGISTRY)/dev-timdex-mario:`git describe --always` \
		-t dev-timdex-mario:latest .

publish: dist ## Build, tag and push
	aws --profile default ecr get-login-password --region us-east-1 \
	| docker login --username AWS --password-stdin $(ECR_REGISTRY)
	docker push $(ECR_REGISTRY)/dev-timdex-mario:latest
	docker push $(ECR_REGISTRY)/dev-timdex-mario:`git describe --always`

promote: ## Promote the current staging build to production
	docker login -u AWS -p $$(aws ecr get-login-password --region us-east-1) $(ECR_REGISTRY)
	docker pull $(ECR_REGISTRY)/mario-stage:latest
	docker tag $(ECR_REGISTRY)/mario-stage:latest $(ECR_REGISTRY)/mario-prod:latest
	docker tag $(ECR_REGISTRY)/mario-stage:latest $(ECR_REGISTRY)/mario-prod:$(DATETIME)
	docker push $(ECR_REGISTRY)/mario-prod:latest
	docker push $(ECR_REGISTRY)/mario-prod:$(DATETIME)
