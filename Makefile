APPNAME := dunai
GO := go
LDFLAGS := -s -w
SUBMAKE := make --no-print-directory
ECR_URL := 193635214029.dkr.ecr.eu-central-1.amazonaws.com/dunai
export GOPATH=$(shell realpath .go)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

init:
	go get

run: ## Run in docker
	docker-compose -f docker-compose-dev.yml up --build
	# go run .

build-prod:
	docker-compose build

run-prod: build-prod
	docker-compose up

dbshell:
	docker-compose exec db psql -U dunai

push: build-prod
	`AWS_PROFILE=adunai aws ecr get-login --no-include-email`
	docker push $(ECR_URL)
	docker logout

test-gemini:
	echo gemini://dun.ai/ | ncat -C --ssl 127.0.0.1 1965 --ssl-servername dun.ai
