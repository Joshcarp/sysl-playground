package files

import "testing"

import "fmt"

func TestFiles(t *testing.T) {
	fs := NewFs("root")
	_, err := fs.Create("tmp.sysl")
	if err != nil {
		panic(err)
	}
	fmt.Println("\n\n\n\nd")
	fs.printall()
	fmt.Println("\n\n\n\nd")
	_, err2 := fs.Open("tmp.sysl")
	if err2 != nil {
		panic(err2)
	}
	// panic(false)

}
