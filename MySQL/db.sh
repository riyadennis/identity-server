#!/bin/sh

docker network create backend

docker build -t mysql-development .

docker run --rm --name=mysql-development \
-e MYSQL_USER='root' -e MYSQL_ROOT_PASSWORD='root' \
-e MYSQL_DATABASE='identity_db' -p=3308:3306 \
--network=backend mysql-development

docker network connect --alias=mysql-development backend mysql-development