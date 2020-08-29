.PHONY: help
all: help
APPNAME = rashomon
PID = $(shell ps -ef | grep ${app-name} | grep -v 'grep' | awk '{print $$2}')

## install: install missing dependencies
install:
	go mod tidy
## build: compile the app
build:
	go build -o bin/${APPNAME} .

## run: run the app
run:
	./bin/${APPNAME}

## stop : stop the app
stop: 
	@kill -15 ${PID}


## test: test the app
test:
	go test -cover ./...

help: Makefile
	@echo
	@echo " Choose a command run in "$(APPNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
