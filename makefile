.PHONY: build run
build:
	go build -o dkhp-uit main.go
run: build
	./dkhp-uit