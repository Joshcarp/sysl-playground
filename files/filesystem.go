package files

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/afero"
)

type FsStruct struct {
	name  string
	files []*File
	dirs  []*Fs
}
type Fs struct {
	fs *FsStruct
}

func NewFs(name string) Fs {
	return Fs{fs: &FsStruct{name: name}}
}

// Create creates a file in the filesystem, returning the file and an
// error, if any happens.
func (fs Fs) Create(name string) (afero.File, error) {
	newFile := newFile(name)
	fs.fs.files = append(fs.fs.files, &newFile)
	fs.printall()
	return newFile, nil
}
func (fs Fs) printall() {
	for i, e := range fs.fs.files {
		fmt.Println(i, e)
	}
}

// Mkdir creates a directory in the filesystem, return an error if any
// happens.
func (fs Fs) Mkdir(name string, perm os.FileMode) error {

	fs.fs.dirs = append(fs.fs.dirs, &Fs{fs: &FsStruct{name: name}})
	return nil
}

// MkdirAll creates a directory path and all parents that does not exist
// yet.
func (fs Fs) MkdirAll(path string, perm os.FileMode) error {
	a := strings.Split(path, "/")
	fs.mkdirAllHelper(a)
	return nil
}

func (fs Fs) mkdirAllHelper(paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	nextDir := &Fs{fs: &FsStruct{name: paths[0]}}
	fs.fs.dirs = append(fs.fs.dirs, nextDir)
	return nextDir.mkdirAllHelper(paths[1:])
}

// Open opens a file, returning it or an error, if any happens.
func (fs Fs) Open(name string) (afero.File, error) {
	for _, file := range fs.fs.files {
		if file.Name() == name {
			return file, nil
		}
	}
	return nil, os.ErrNotExist
}

// OpenFile opens a file using the given flags and the given mode.
func (fs Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	panic("OpenFile not implemented") // TODO: Implement

}

// Remove removes a file identified by name, returning an error, if any
// happens.
func (fs Fs) Remove(name string) error {
	for i, file := range fs.fs.files {
		if file.Name() == name {
			fs.fs.files = append(fs.fs.files[:i], fs.fs.files[i+1])
			return nil
		}
	}
	return nil
}

// RemoveAll removes a directory path and any children it contains. It
// does not fail if the path does not exist (return nil).
func (fs Fs) RemoveAll(path string) error {
	for i, dir := range fs.fs.dirs {
		if dir.fs.name == path {
			fs.fs.dirs = append(fs.fs.dirs[:i], fs.fs.dirs[i+1])
			return nil
		}
	}
	return nil
}

// Rename renames a file.
func (fs Fs) Rename(oldname string, newname string) error {
	for _, file := range fs.fs.files {
		if file.Name() == oldname {
			file.f.info = file.f.info.setName(newname)
		}
	}
	return nil
}

// Stat returns a FileInfo describing the named file, or an error, if any
// happens.
func (fs Fs) Stat(name string) (os.FileInfo, error) {
	file, err := fs.Open(name) // TODO: Implement
	if err != nil {
		return nil, err
	}
	fInfo, err := file.Stat()

	return fInfo, err
}

// The name of this FileSystem
func (fs Fs) Name() string {
	return fs.fs.name
}

//Chmod changes the mode of the named file to mode.
func (fs Fs) Chmod(name string, mode os.FileMode) error {
	panic("Chmod not implemented") // TODO: Implement
}

//Chtimes changes the access and modification times of the named file
func (fs Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	panic("Chtimes not implemented") // TODO: Implement
}
