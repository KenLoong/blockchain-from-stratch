build:
	go build -o ./bin/warson-project

run: build
	./bin/warson-project
test:
	go test -v ./...