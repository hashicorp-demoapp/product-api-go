CONTAINER_NAME=hashicorpdemoapp/product-api
DB_CONTAINER_NAME=hashicorpdemoapp/product-api-db
CONTAINER_VERSION=v0.0.17

test_functional:
	shipyard run ./blueprint
	cd ./functional_tests && go test -v -run.test true ./..
	shipyard destroy

build_db:
	docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	docker buildx create --name multi-db || true
	docker buildx use multi-db
	docker buildx inspect --bootstrap
	docker buildx build --platform linux/arm64,linux/amd64 \
		-t ${DB_CONTAINER_NAME}:${CONTAINER_VERSION} \
		-f ./database/Dockerfile \
		./database \
		--push
	docker buildx rm multi-db

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/amd64/product-api

build_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/arm64/product-api

build_docker: build_linux build_arm64
	docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	docker buildx create --name multi || true
	docker buildx use multi
	docker buildx inspect --bootstrap
	docker buildx build --platform linux/arm64,linux/amd64 \
		-t ${CONTAINER_NAME}:${CONTAINER_VERSION} \
		-f ./Dockerfile \
		./bin \
		--push
	docker buildx rm multi