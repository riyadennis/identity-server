# identity-server [![CircleCI](https://circleci.com/gh/riyadennis/identity-server.svg?style=svg)](https://circleci.com/gh/riyadennis/identity-server) [![Go Reference](https://pkg.go.dev/badge/github.com/riyadennis/identity-server.svg)](https://pkg.go.dev/github.com/riyadennis/identity-server) <a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-33%25-brightgreen.svg?longCache=true&style=flat)</a>

To Run the service locally we need .env file set with the following values:

```
PORT="8089"
ENV="dev"
MYSQL_DATABASE="identity"
MYSQL_USERNAME="username"
MYSQL_PASSWORD="password"
MYSQL_PORT="3306"
MYSQL_HOST="127.0.0.1"
BASE_PATH="/"
```

For the successful running of tests we need an .env_test file with following values:

```
PORT="8089"
ENV="test"
MYSQL_DATABASE="test"
MYSQL_USERNAME="newuser"
MYSQL_PASSWORD="password"
MYSQL_PORT="3306"
MYSQL_HOST="127.0.0.1"
BASE_PATH="/"
```
