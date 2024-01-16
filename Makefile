all:
	go build -mod=vendor -o dist/webpty.bin main.go

test:
	go test -v -coverprofile dist/cover.out ./...
	go tool cover -html dist/cover.out -o dist/cover.html

tunnel:
	docker build -f ./webfleet/Dockerfile -t machines/webpty-tunnel:latest .
	docker push machines/webpty-tunnel:latest
