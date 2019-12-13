# Derived from official mysql image (our base image)
FROM mysql
ENV MYSQL_DATABASE identity-server
ENV MYSQL_USER identity
ENV MYSQL_PASSWORD identity
COPY . /docker-entrypoint-initdb.d/