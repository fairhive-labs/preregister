build: clean
	go build -o bin/preregister-api -v ./cmd/api/*
run: clean build
	./bin/preregister-api
clean:
	rm -rf ./bin
test:
	go test -v ./...