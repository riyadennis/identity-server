docker-build:
	docker build -t riyadennis/identity-server .

docker-push:
	# need to do the push with a new tag
	docker push riyadennis/identity-server:1.0.2

helm-install:
	helm install identity ./zarf/identity

helm-uninstall:
	helm uninstall identity

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
	export GITLAB_TOKEN="your token"
