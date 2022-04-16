build: clean
	go build -o bin/preregister-api -v ./cmd/api/main.go
run: clean build
	./bin/preregister-api
clean:
	rm -rf ./bin
test:
	go clean -testcache
	go test -v ./... 