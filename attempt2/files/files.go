package files

import "os"

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
	p = f.f.contents
	return 0, nil
}

func (f File) ReadAt(p []byte, off int64) (n int, err error) {
	panic("ReadAt not implemented") // TODO: Implement
}

func (f File) Seek(offset int64, whence int) (int64, error) {
	panic("Seek not implemented") // TODO: Implement
}

func (f File) Write(p []byte) (n int, err error) {
	f.f.contents = p
	return 0, nil
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
