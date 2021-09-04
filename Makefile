minikube-start:
	minikube start --vm-driver=virtualbox --disk-size=30g

docker-build:
	docker build -t riyadennis/identity-server:1.3.0 .

docker-push:
	# need to do the push with a new tag
	docker push riyadennis/identity-server:1.3.0

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
	git tag -a v0.1.4 -m "fix for duplicate entries"
	git push origin v0.1.4

githubToken:
	export GITHUB_TOKEN="your token"

exportEnv:
	export $(xargs < .env)
