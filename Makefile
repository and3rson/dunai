ECS_URL = 193635214029.dkr.ecr.eu-central-1.amazonaws.com/dunai

dev: | build
	docker-compose -f ops/dev/docker-compose.yml -p dunai up

prod: | build
	docker-compose -f ops/prod/docker-compose.yml -p dunai up

bash:
	docker exec -it dunai_app_1 bash

build:
	docker build -t dunai .

run: | build
	docker run --rm --name dunai -it dunai

deploy: | build
	docker tag dunai ${ECS_URL}:latest
	docker push ${ECS_URL}:latest

