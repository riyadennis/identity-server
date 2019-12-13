docker-build:
	docker-compose up --build
docker-run:
	docker run --rm -p 8080:8080  identity-server
db-build:
	docker build -t identity-db db/
db-run:
	docker run -p 3306:3306 --name identity-db -e MYSQL_ROOT_PASSWORD=root -d mysql
db-login:
	docker exec -it identity-db bash