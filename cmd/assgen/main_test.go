package main

import (
	"os"
	"testing"
)

func TestGenerateManifest(t *testing.T) {
	if exists("../../test/dist") {
		os.RemoveAll("../../test/dist")
	}

	ret := GenerateManifest("../../test/config.yaml", "./dist")
	if ret != 0 {
		t.Fatalf("GenerateManifest returned %d", ret)
	}

	if !exists("../../test/dist") {
		t.Fatal("dist directory not created")
	}

	testFiles := []string{
		"../../test/dist/test.js",
		"../../test/dist/other.js",
		"../../test/dist/test.css",
	}

	for _, f := range testFiles {
		if !exists(f) {
			t.Fatalf("test file not created: %s", f)
		}
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
