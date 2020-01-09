package files

import (
	"io"
	"os"
)

type FileStruct struct {
	info     fileInfo
	contents []byte
}

func (f FileStruct) String() string {
	return f.info.Name()
}

type File struct {
	f *FileStruct
}

func newFileStruct() FileStruct {
	return FileStruct{}
}
func newFile(name string) File {
	info := newFileInfo(name)
	f := FileStruct{info: info}
	return File{f: &f}
}

func (f File) Close() error {
	return nil
}

func (f File) Read(p []byte) (n int, err error) {

	lenp := cap(p)
	lenf := len(f.f.contents) - int(f.f.info.offset)
	if lenf == 0 {
		f.f.info.offset = 0
		return 0, io.EOF
	}
	if lenp > lenf {
		// p = f.f.contents
		n = copy(p, f.f.contents[f.f.info.offset:lenf])
		f.f.info.offset = int64(lenf)
		return n, nil
	}
	f.f.info.offset = int64(lenp)
	// p = f.f.contents[:lenp]
	n = copy(p, f.f.contents[f.f.info.offset:lenp])
	return n, nil
}

func (f File) ReadAt(p []byte, off int64) (n int, err error) {

	// n, err := f.content.ReadAt(b, off)

	panic("not implemented ")
}

func (f File) Seek(offset int64, whence int) (int64, error) {
	panic("Seek not implemented") // TODO: Implement
}

func (f File) Write(p []byte) (int, error) {
	n := copy(f.f.contents, p)
	return n, nil
}

func (f File) WriteAt(p []byte, off int64) (n int, err error) {
	panic("WriteAt not implemented") // TODO: Implement
}

func (f File) Name() string {
	return f.f.info.Name()
}

func (f File) Readdir(count int) ([]os.FileInfo, error) {
	panic("Readdir not implemented") // TODO: Implement
}

func (f File) Readdirnames(n int) ([]string, error) {
	panic("Readdirnames not implemented") // TODO: Implement
}

func (f File) Stat() (os.FileInfo, error) {

	return f.f.info, nil
	// panic("Stat not implemented") // TODO: Implement
}

func (f File) Sync() error {
	panic("Sync not implemented") // TODO: Implement
}

func (f File) Truncate(size int64) error {
	panic("Truncate not implemented") // TODO: Implement
}

func (f File) WriteString(s string) (ret int, err error) {
	panic("WriteString not implemented") // TODO: Implement
}
