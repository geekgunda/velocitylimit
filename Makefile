all: build run

build:
		cd ${CURDIR}/cmd/app; go get -d; go clean -r; go build;

run:
		cd ${CURDIR}/cmd/app/; ./app

test:
		go test -v
