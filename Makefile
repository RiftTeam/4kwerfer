NAME=4kwerfer

default: build

build: deps
	go build -o $(NAME) .


deps:
	go get github.com/pointlander/peg
	git submodule update --init --recursive
