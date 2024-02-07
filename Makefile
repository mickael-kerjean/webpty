all:
	go build -mod=vendor -o dist/webpty.bin main.go

install:
	mv dist/webpty.bin /usr/local/bin/webpty.bin
	mv systemctl /etc/systemd/system/webpty.service

test:
	go test -v -coverprofile dist/cover.out ./...
	go tool cover -html dist/cover.out -o dist/cover.html

tunnel:
	docker build -f ./webfleet/Dockerfile -t machines/webpty-tunnel:latest .
	docker push machines/webpty-tunnel:latest
