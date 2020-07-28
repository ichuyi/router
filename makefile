run:
	GOARCH=amd64 GOOS=linux go build -o server server.go
deploy:
	rm server_*
	docker build -t demo:test .