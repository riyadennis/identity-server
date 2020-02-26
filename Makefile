docker-build:
	docker build -t riyadennis/identity-server .

docker-run:
	docker run --rm -p 8095:8081  riyadennis/identity-server

docker-push:
	docker push riyadennis/identity-server:latest