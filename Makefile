lint:
	golint .

build:
	go build -o bin/bg_manager

run: build
	./bin/bg_manager

test:
	go test -v ./...

install: build
	sudo cp ./bin/bg_manager /usr/bin/