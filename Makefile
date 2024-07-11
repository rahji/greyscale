BINARY_NAME=greyscale
 
all: build
 
build:
	go build -o ./bin/${BINARY_NAME} main.go
 
test:
	go test -v main.go
 
run:
	go build -o ./bin/${BINARY_NAME} main.go
	./bin/${BINARY_NAME}
 
clean:
	go clean
	rm ./bin/${BINARY_NAME}
