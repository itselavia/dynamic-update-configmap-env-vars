build-linux:
	GOOS=linux go build -o bin/update_env

docker-build: build-linux
	docker build -t itselavia/dynamic-update-env .

docker-push: docker-build
	docker push itselavia/dynamic-update-env
	rm -rf bin/