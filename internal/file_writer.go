package internal

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func NewFileWriter(input io.Reader, path string, fileType string) (FileWriter, error) {
	if !strings.HasPrefix(fileType, ".") {
		return FileWriter{}, fmt.Errorf("invalid file type: %s", fileType)
	}

	return FileWriter{
		input:    input,
		path:     path,
		fileType: fileType,
	}, nil
}

type FileWriter struct {
	input    io.Reader
	path     string
	fileType string
}

func (writer *FileWriter) Run(noop bool) error {
	file, err := writer.getCopyDestination(noop)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, writer.input)
	if err != nil {
		return err
	}

	return nil
}

func (writer *FileWriter) getCopyDestination(noop bool) (io.WriteCloser, error) {
	if noop {
		return NopWriteCloser{}, nil
	}

	targetDir := filepath.Dir(writer.path)
	_, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		log.Println("output path does not exist, creating...")
		err := os.Mkdir(targetDir, 0755)
		if err != nil {
			return nil, err
		}
		log.Printf("output path created: %s", targetDir)
	} else {
		log.Printf("output path exists: %s", targetDir)
	}

	return os.Create(fmt.Sprintf("%s%s", writer.path, writer.fileType))
}

type NopWriteCloser struct{}

func (NopWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (NopWriteCloser) Close() error {
	return nil
}
