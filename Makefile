build: clean
	go build -o bin/preregister-api -v ./cmd/api/*
run: build
	./bin/preregister-api
clean:
	rm -rf ./bin
test:
	go test -v ./...