package internal

import (
	"bytes"
	"os"
	"testing"
)

func TestFileWriter(t *testing.T) {
	reader := bytes.NewReader([]byte("hello, world!"))
	_, err := NewFileWriter(reader, "../test/dist/test", ".txt")
	if err != nil {
		t.Fatalf("error creating file writer: %v", err)
	}
}

func TestFileWriter_Run(t *testing.T) {
	reader := bytes.NewReader([]byte("hello, world!"))
	fileWriter, err := NewFileWriter(reader, "../test/dist/test", ".txt")
	if err != nil {
		t.Fatalf("error creating file writer: %v", err)
	}

	err = fileWriter.Run()
	if err != nil {
		t.Fatalf("error running the FileWriter: %v", err)
	}

	if !exists("../test/dist") {
		t.Fatalf("output dir not created")
	}

	if !exists("../test/dist/test.txt") {
		t.Fatalf("output file not created")
	}

	content, err := os.ReadFile("../test/dist/test.txt")
	if err != nil {
		t.Fatalf("error reading output file: %v", err)
	}

	if "hello, world!" != string(content) {
		t.Fatalf("output content incorrect: '%s', expected: 'hello, world!'", content)
	}

	os.RemoveAll("../test/dist")
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
