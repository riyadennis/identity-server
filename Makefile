minikube-start:
	minikube start --vm-driver=virtualbox --disk-size=30g

docker-build:
	docker build -t riyadennis/identity-server:1.4.0 .

docker-push:
	# need to do the push with a new tag
	docker push riyadennis/identity-server:1.4.0

helm-install:
	helm install identity ./zarf/identity

helm-uninstall:
	helm uninstall identity

#helps to fetch service URL to access the pod
minikube-services:
	minikube service list

service-url:
	minikube service identity --url

mysql-install:
	helm install my-sql -f mysql-chart/values.yaml bitnami/mysql

docker-run:
	docker run -d --rm -e ENV='prod' -e PORT=":8095" -e KEY="this-should-be-secret-shared-only-to-client" \
    -e ISSUER="riya-dennis" -p 8095:8095 --name identity-server --network backend riyadennis/identity-server

docker-run-test:
	docker run -d --rm -e ENV='test' -e PORT=":8095" -e KEY="this-should-be-secret-shared-only-to-client" \
    -e ISSUER="riya-dennis" -p 8088:8088 --name identity-server-test --network backend riyadennis/identity-server

tag:
	git tag -a v0.1.7 -m "fixed docker build and push"
	git push origin v0.1.7

githubToken:
	export GITHUB_TOKEN="your token"

exportEnv:
	export $(xargs < .env)
    GOBIN ?= $$(go env GOPATH)/bin

install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

check-coverage: install-go-test-coverage
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=testcoverage.yaml

docs-fmt:
	swag fmt

docs-generate:
	swag init -g app/auth-api/main.go

claude:
	npx @anthropic-ai/claude-code chat

docker-compose-up:
	docker-compose build --no-cache
	docker-compose up
install-proto:
    brew install protobuf
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
check-proto:
    protoc --version
generate-proto:
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative app/auth-api/gRPC/identity.proto 