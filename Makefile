CONTAINER_NAME=hashicorpdemoapp/product-api
DB_CONTAINER_NAME=hashicorpdemoapp/product-api-db
CONTAINER_VERSION=v0.0.10

# Run tests continually with  a watcher
autotest:
	filewatcher --idle-timeout 24h -x **/functional_tests gotestsum --format standard-verbose

build_db:
	cd database && docker build -t ${DB_CONTAINER_NAME}:${CONTAINER_VERSION} .

build_linux:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/product-api

build_docker: build_linux
	docker build -t ${CONTAINER_NAME}:${CONTAINER_VERSION} .

push_docker: build_docker
	docker push ${CONTAINER_NAME}:${CONTAINER_VERSION}
