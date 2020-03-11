#!/bin/sh

source .env
docker network create backend

docker build -t mysql-development .

docker run --rm --name=mysql-development \
-e MYSQL_USER='{$MYSQL_USER}' -e MYSQL_ROOT_PASSWORD='{$MYSQL_ROOT_PASSWORD}' \
-e MYSQL_DATABASE='{$MYSQL_DATABASE}' -p=3308:3306 \
--network=backend mysql-development

docker network connect --alias=mysql-development backend mysql-development