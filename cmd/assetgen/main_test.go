package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateManifest(t *testing.T) {
	if exists("../../test/dist") {
		os.RemoveAll("../../test/dist")
	}

	configPath := "../../test/config.yaml"
	configPath, _ = filepath.Abs(configPath)
	appCtx := AppContext{configPath: configPath, configDir: filepath.Dir(configPath), outPath: ""}
	ret := GenerateManifest(appCtx)
	if ret != 0 {
		t.Fatalf("GenerateManifest returned %d", ret)
	}

	if !exists("../../test/dist") {
		t.Fatal("dist directory not created")
	}

	if !exists("../../test/dist/subdir") {
		t.Fatal("dist/subdir directory not created")
	}

	testFiles := []string{
		"../../test/dist/test.js",
		"../../test/dist/other.js",
		"../../test/dist/test.css",
		"../../test/dist/subdir/inner.js",
		"../../test/dist/random.txt",
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
