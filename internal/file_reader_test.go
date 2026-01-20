package internal

import (
	"testing"
)

func TestNewFileReader_Exists(t *testing.T) {
	path := "../test/test.js"

	fileReader, err := NewFileReader(path)
	if err != nil {
		t.Fatalf("file reader err: %v", err)
	}

	if fileReader.Type() != ".js" {
		t.Fatalf("incorrect type: %v, expected '.js'", fileReader.Type())
	}
}

func TestNewFileReader_NoFile(t *testing.T) {
	path := "../test/nonexistent.js"

	_, err := NewFileReader(path)
	if err == nil {
		t.Fatalf("file reader did not return an error\n")
	}
}
