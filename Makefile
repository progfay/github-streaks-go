all: ghs.out

ghs.out: ./main.go
	go build -o ghs.out .
