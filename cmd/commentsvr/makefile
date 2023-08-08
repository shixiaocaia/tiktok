.PHONY: all
TARGET := commentsvr
GOENV := GOOS=linux GOARCH=amd64
GOMACENV := GOOS=darwin GOARCH=amd64
all:
	CGO_ENABLED=0 ${GOENV} go build -o ./bin/${TARGET}

clean:
	rm -rf ${TARGET}
format:
	gofmt -w .
	goimports -w .
	golint ./...
test:
	go test --cover -gcflags=-l ./...

build:
	CGO_ENABLED=0 ${GOENV} go build -o ./bin/${TARGET}