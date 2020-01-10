all: main

main: *.go
	GOOS=js GOARCH=wasm go1.12.4 build -mod vendor -o main.wasm main.go
serve:
	go run server/server.go