package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Styles  []string `yaml:"styles"`
	Scripts []string `yaml:"scripts"`
	Out     *string  `yaml:"out"`
}

type Asset struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

type Manifest struct {
	Styles  []Asset `json:"styles"`
	Scripts []Asset `json:"scripts"`
}

func GenerateManifest(config string, out string) int {
	configFilePath, err := filepath.Abs(config)
	if err != nil {
		log.Println(err)
		return 1
	}

	configFileDir := filepath.Dir(configFilePath)
	outputPath := filepath.Join(configFileDir, out)

	log.Println("ensuring output path")
	_, err = os.Stat(outputPath)
	if os.IsNotExist(err) {
		log.Println("output path does not exist, creating...")
		err := os.Mkdir(outputPath, 0755)
		if err != nil {
			log.Println(err)
			return 1
		}
		log.Println("output path created")
	} else {
		log.Println("output path exists")
	}

	configDoc, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Println(err)
		return 1
	}

	c := Config{}

	err = yaml.Unmarshal(configDoc, &c)
	if err != nil {
		log.Println(err)
		return 1
	}

	scripts := c.Scripts
	styles := c.Styles

	manifest := Manifest{}

	log.Println("processing scripts")
	scriptAssets, err := processGlobs(scripts, configFileDir, outputPath)
	manifest.Scripts = scriptAssets

	log.Println("processing styles")
	styleAssets, err := processGlobs(styles, configFileDir, outputPath)
	manifest.Styles = styleAssets

	manifestContent, err := json.Marshal(manifest)
	if err != nil {
		log.Println(err)
		return 1
	}

	manifestPath := filepath.Join(outputPath, "manifest.json")
	err = writeManifest(manifestPath, manifestContent)
	if err != nil {
		log.Println(err)
		return 1
	}

	log.Println("manifest written")

	return 0
}

func processGlobs(globs []string, configFileDir string, outputPath string) ([]Asset, error) {
	results := make([]Asset, 0)

	for _, script := range globs {
		fullpath := filepath.Join(configFileDir, script)
		matches, err := filepath.Glob(fullpath)
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			rel, err := filepath.Rel(configFileDir, match)
			if err != nil {
				return nil, err
			}

			destPath := filepath.Join(outputPath, rel)
			log.Printf("copying %s\n", rel)
			hash, err := copyFile(match, destPath)
			if err != nil {
				return nil, err
			}

			results = append(results, Asset{Path: rel, Hash: hash})
		}
	}

	return results, nil
}

func copyFile(from string, to string) (string, error) {
	src, err := os.Open(from)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(to)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	hasher := sha256.New()
	data := io.TeeReader(src, hasher)

	_, err = io.Copy(dst, data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func writeManifest(path string, content []byte) error {
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.Write(content)
	return err
}
