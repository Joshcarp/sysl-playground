main: *.go
	GOARCH=wasm GOOS=js go build -o main.wasm main.go
