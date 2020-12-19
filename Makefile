all: ghs.out

ghs.out: src/main.go
	go build -o ghs.out ./src
