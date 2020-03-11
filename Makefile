docker-build:
	docker build -t riyadennis/identity-server .

docker-run:
	docker run --rm -e ENV='prod' -e CONFIG_FILE='/identity-server/etc/config.yaml' \
    -p 8095:8095 --name identity-server --network backend riyadennis/identity-server

docker-run-test:
	docker run --rm  -e ENV='test' -e CONFIG_FILE='/identity-server/etc/config_test.yaml' \
    -p 8088:8088 --name identity-server-test riyadennis/identity-server

docker-push:
	docker push riyadennis/identity-server:latest