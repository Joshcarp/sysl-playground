package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewMemMapFs()
	fs.Create("/tmp.sysl")

	fmt.Println(filepath.Abs("/tmp"))
	this := afero.NewBasePathFs(fs, "/").(*afero.BasePathFs)

	currentPath, err := this.RealPath("tmp")
	fmt.Println("world", currentPath, err)
}
