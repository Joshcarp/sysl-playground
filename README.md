# sysl-playground
sysl-playground written in go vecty [https://github.com/gopherjs/vecty]

## Requirements
 - Go 1.12.4 or go 1.13 compiled from source with `initialSize` increased to accomodate the large binary size of sysl
    - See [https://github.com/golang/go/issues/34395] for more details
- Playground uses sysl package from vendor package which has a modified version of sysl in it; Don't update vendor package

## Features
- Sysl diagram generation from sysl code
- Share Diagrams and code via url query parameters; NOTE: Can only share sysl files with only ~2600 characters due to url length limit

