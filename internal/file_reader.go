package internal

import (
	"io"
	"os"
	"path/filepath"
)

func NewFileReader(path string) (FileReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return FileReader{}, err
	}

	extension := filepath.Ext(path)

	return FileReader{file: file, fileType: extension}, nil
}

type FileReader struct {
	file     *os.File
	fileType string
}

func (f *FileReader) Reader() io.Reader {
	return f.file
}

func (f *FileReader) Closer() io.Closer {
	return f.file
}

func (f *FileReader) Type() string {
	return f.fileType
}
