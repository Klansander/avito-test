all: sfera

sfera:
	swag init -g ./cmd/app/main.go -o ./docs
	GOOS=linux GOARCH=amd64 go build -o cmd/app/sfera-server1 ./cmd/app/main.go

swag:	
	swag init -g ./cmd/app/main.go -o ./docs

clean:
	rm $(SFERA)

rebuild: clean all
