docker-build:
	docker build -t identity-server .
docker-run:
	docker run --rm -p 8081:8081  identity-server