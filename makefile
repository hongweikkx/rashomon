app-name = rashomon
all : build run
build:
	go build -o bin/${app-name} .

run:
	./bin/${app-name}

test:
	go test -cover ./...
