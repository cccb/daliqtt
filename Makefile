

all:
	go get .
	go build

linux:
	GOARCH=amd64 GOOS=linux go build 

