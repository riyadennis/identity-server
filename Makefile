docker-build:
	docker-compose up --build
docker-run:
	docker run --rm -p 8080:8080  identity-server