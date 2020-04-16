docker-build:
	docker build -t riyadennis/identity-server .

docker-run:
	docker run -d --rm -e ENV='prod' -e PORT=":8095" -e KEY="this-should-be-secret-shared-only-to-client" \
    -e ISSUER="riya-dennis"  -e MYSQL_USER="root" -e MYSQL_ROOT_PASSWORD="root" -e MYSQL_DATABASE="identity_db" \
    -e MYSQL_PORT="3306" -e MYSQL_HOST="mysql-development" \
    -p 8095:8095 --name identity-server --network backend riyadennis/identity-server

docker-run-test:
	docker run -d --rm -e ENV='test' -e PORT=":8095" -e KEY="this-should-be-secret-shared-only-to-client" \
    -e ISSUER="riya-dennis" -p 8088:8088 --name identity-server-test --network backend riyadennis/identity-server

docker-push:
	docker push riyadennis/identity-server:latest