HOSTNAME=registry.terraform.io
NAMESPACE=CruGlobal
NAME=redisdb
BINARY=terraform-provider-${NAME}
VERSION=0.1
OS_ARCH=darwin_amd64

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	./redis.sh update
	./redis.sh start
	./redis.sh test
	./redis.sh stop

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
