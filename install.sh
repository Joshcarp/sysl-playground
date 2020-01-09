#!/bin/bash

GOOS=js GOARCH=wasm go1.12.4 build -mod vendor -o main.wasm markdown.go

goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'