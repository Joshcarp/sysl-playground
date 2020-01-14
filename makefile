all: main

main: *.go
	GOOS=js GOARCH=wasm go1.12.4 build -mod vendor -o static/main.wasm main.go
serve: main
	go run server/server.go