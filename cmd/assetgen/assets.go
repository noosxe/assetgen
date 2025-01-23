package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

type AppContext struct {
	configPath string
	configDir  string
	outPath    string
	noCopy     bool
	noManifest bool
}

type Config struct {
	Styles  []string `yaml:"styles"`
	Scripts []string `yaml:"scripts"`
	Random  []string `yaml:"random"`
	Out     *string  `yaml:"out"`
}

type Asset struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

type Manifest struct {
	Styles  []Asset `json:"styles"`
	Scripts []Asset `json:"scripts"`
	Random  []Asset `json:"random"`
}

func GenerateManifest(appCtx AppContext) int {
	c, err := ReadConfig(appCtx.configPath)
	if err != nil {
		log.Println(err)
		return 1
	}

	if appCtx.outPath == "" {
		if c.Out == nil {
			log.Println("no output path specified")
			return 1
		}

		out := *c.Out
		if filepath.IsAbs(out) {
			appCtx.outPath = out
		} else {
			appCtx.outPath = filepath.Join(appCtx.configDir, out)
		}
	}

	log.Printf("output path is %s", appCtx.outPath)
	log.Println("ensuring output path")
	_, err = os.Stat(appCtx.outPath)
	if os.IsNotExist(err) {
		log.Println("output path does not exist, creating...")
		err := os.Mkdir(appCtx.outPath, 0755)
		if err != nil {
			log.Println(err, appCtx.outPath)
			return 1
		}
		log.Println("output path created")
	} else {
		log.Println("output path exists")
	}

	scripts := c.Scripts
	styles := c.Styles
	random := c.Random

	manifest := Manifest{}

	log.Println("processing scripts")
	scriptAssets, err := processGlobs(appCtx, scripts, appCtx.configDir, appCtx.outPath)
	if err != nil {
		log.Println(err)
		return 1
	}

	manifest.Scripts = scriptAssets

	log.Println("processing styles")
	styleAssets, err := processGlobs(appCtx, styles, appCtx.configDir, appCtx.outPath)
	if err != nil {
		log.Println(err)
		return 1
	}

	manifest.Styles = styleAssets

	log.Println("processing random assets")
	randomAssets, err := processGlobs(appCtx, random, appCtx.configDir, appCtx.outPath)
	if err != nil {
		log.Println(err)
		return 1
	}

	manifest.Random = randomAssets

	manifestContent, err := json.Marshal(manifest)
	if err != nil {
		log.Println(err)
		return 1
	}

	manifestPath := filepath.Join(appCtx.outPath, "manifest.json")
	err = writeManifest(manifestPath, manifestContent)
	if err != nil {
		log.Println(err)
		return 1
	}

	log.Println("manifest written")

	return 0
}

func ReadConfig(path string) (*Config, error) {
	configDoc, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := Config{}

	err = yaml.Unmarshal(configDoc, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func processGlobs(appCtx AppContext, globs []string, configFileDir string, outputPath string) ([]Asset, error) {
	results := make([]Asset, 0)

	for _, script := range globs {
		fullpath := filepath.Join(configFileDir, script)
		basepath, pattern := doublestar.SplitPattern(fullpath)

		fsys := os.DirFS(basepath)
		matches, err := doublestar.Glob(fsys, pattern)

		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			glued := filepath.Join(basepath, match)
			rel, err := filepath.Rel(configFileDir, glued)
			if err != nil {
				return nil, err
			}

			destPath := filepath.Join(outputPath, rel)
			log.Printf("copying %s\n", rel)
			hash, err := copyFile(appCtx, glued, destPath)
			if err != nil {
				return nil, err
			}

			results = append(results, Asset{Path: rel, Hash: hash})
		}
	}

	return results, nil
}

func copyFile(appCtx AppContext, from string, to string) (string, error) {
	src, err := os.Open(from)
	if err != nil {
		return "", err
	}
	defer src.Close()

	if !appCtx.noCopy {
		destDir := filepath.Dir(to)
		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			return "", err
		}
	}

	dst, err := getCopyDestination(appCtx, to)
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

func getCopyDestination(appCtx AppContext, to string) (io.WriteCloser, error) {
	if appCtx.noCopy {
		return NopWriteCloser{}, nil
	}

	return os.Create(to)
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

type NopWriteCloser struct{}

func (NopWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (NopWriteCloser) Close() error {
	return nil
}
