package internal

import (
	"encoding/json"
	"os"
)

type Asset struct {
	Id        *string `json:"id,omitempty"`
	MediaType string  `json:"mediaType"`
	Path      string  `json:"path"`
	Hash      string  `json:"hash"`
	Preload   bool    `json:"preload,omitempty"`
}

type Manifest struct {
	Assets []Asset `json:"assets"`
}

func WriteManifest(path string, manifest *Manifest) error {
	manifestContent, err := json.Marshal(manifest)
	if err != nil {
		return err
	}

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.Write(manifestContent)
	return err
}
