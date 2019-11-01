NAME=wgrest
DESTDIR?=./dist

linux-amd64:
	GOOS=linux ARCH=amd64 go build -o $(DESTDIR)/$(NAME)-linux-amd64 cmd/wgrest-server/main.go

all: linux-amd64