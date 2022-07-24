.PHONY:
.SILENT:

build:
	go build -o ./.bin/event cmd/event/main.go
run: build
	sudo ./.bin/event