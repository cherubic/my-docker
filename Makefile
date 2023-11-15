

start-dev-env:
	docker build -t dev-env -f Dockerfile.dev .
	docker run -d --name my-docker -v $(PWD):/app dev-env tail -f /dev/null

clean-dev-env:
	if [ -n "$(docker ps -a | grep my-docker)" ]; then docker stop my-docker && docker rm my-docker; fi

.PHONY: start-dev-env clean-dev-env