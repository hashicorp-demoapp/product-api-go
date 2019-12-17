CONTAINER_NAME=hashicorpdemoapp/product-api
CONTAINER_VERSION=latest

# Run tests continually with  a watcher
autotest:
	filewatcher --idle-timeout 24h -x **/functional_tests gotestsum --format standard-verbose

build_linux:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/product-api

build_docker: build_linux
	docker build -t ${CONTAINER_NAME}:${CONTAINER_VERSION} .