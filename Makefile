docker-build:
	docker build -t identity-server .
docker-run:
	docker run --rm -p 8080:8080  identity-server