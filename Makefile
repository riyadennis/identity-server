docker-build:
	docker build -t riyadennis/identity-server .

docker-run:
	docker run --rm  riyadennis/identity-server

docker-push:
	docker push riyadennis/identity-server:latest