NAME=4kwerfer

default: build

build: deps strings
	go build -o $(NAME) .

deps:
	go mod tidy
	go mod download

strings: tools
	# stringer -type=UniformType gl/types.go

tools:
	go get golang.org/x/tools/cmd/stringer
