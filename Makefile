run: clean build
	./bin/api
build: clean
	go build -o bin/api -v ./cmd/api/*.go
clean:
	rm -rf ./bin
test:
	go clean -testcache
	go test ./... 