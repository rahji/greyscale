BINARY_NAME=greyscale
 
all: build
 
build:
	mkdir -p ./bin
	go build -ldflags="-X 'github.com/rahji/greyscale/cmd.version=1.0.2'" -o ./bin/${BINARY_NAME} main.go
test:
	go test -v main.go
 
run:
	mkdir -p ./bin
	go build -ldflags="-X 'github.com/rahji/greyscale/cmd.version=1.0.2'" -o ./bin/${BINARY_NAME} main.go
	./bin/${BINARY_NAME}
 
clean:
	go clean
	rm ./bin/${BINARY_NAME}
