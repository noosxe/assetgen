package main

import (
	"encoding/json"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"

	"github.com/noosxe/assetgen/internal"
)

type AppContext struct {
	configPath string
	configDir  string
	outPath    string
	noCopy     bool
	noManifest bool
}

type Descriptor struct {
	Id      *string `yaml:"id"`
	Glob    string  `yaml:"glob"`
	Preload bool    `yaml:"preload"`
}

type Asset struct {
	Id      *string `json:"id,omitempty"`
	Path    string  `json:"path"`
	Hash    string  `json:"hash"`
	Preload bool    `json:"preload,omitempty"`
}

type Manifest struct {
	Styles  []Asset `json:"styles"`
	Scripts []Asset `json:"scripts"`
	Random  []Asset `json:"random"`
}

func GenerateManifest(appCtx AppContext) int {
	c, err := internal.ReadConfig(appCtx.configPath)
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

	_, err = processGlobs(c.Assets, appCtx.configDir, appCtx.outPath, appCtx)

	manifest := Manifest{}
	manifestContent, err := json.Marshal(manifest)
	if err != nil {
		log.Println(err)
		return 1
	}

	if !appCtx.noManifest {
		manifestPath := filepath.Join(appCtx.outPath, "manifest.json")
		err = writeManifest(manifestPath, manifestContent)
		if err != nil {
			log.Println(err)
			return 1
		}

		log.Println("manifest written")
	} else {
		log.Println("manifest skipped")
	}

	return 0
}

func processGlobs(globs []internal.Input, configFileDir string, outputPath string, appCtx AppContext) ([]internal.Asset, error) {
	results := make([]internal.Asset, 0)

	for _, input := range globs {
		fullpath := filepath.Join(configFileDir, input.Glob)
		basepath, pattern := doublestar.SplitPattern(fullpath)

		fsys := os.DirFS(basepath)
		matches, err := doublestar.Glob(fsys, pattern)

		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			inPath := filepath.Join(basepath, match)
			rel, err := filepath.Rel(configFileDir, inPath)
			if err != nil {
				return nil, err
			}

			outPath := noExt(filepath.Join(outputPath, rel))
			log.Printf("copying %s\n", rel)
			asset := internal.Asset{Id: input.Id, Path: rel, Preload: input.Preload}
			err = pipeline(&asset, input, inPath, outPath, appCtx)
			if err != nil {
				return nil, err
			}

			results = append(results, asset)
		}
	}

	return results, nil
}

func pipeline(asset *internal.Asset, input internal.Input, inPath string, outPath string, appCtx AppContext) error {
	fileReader, err := internal.NewFileReader(inPath)
	if err != nil {
		return err
	}

	lastReader := fileReader.Reader()

	if input.Minify {
		mediaType := mime.TypeByExtension(fileReader.Type())
		minifier := internal.NewMinifier(lastReader, mediaType)
		lastReader = minifier.Reader()
		asset.MediaType = mediaType
	}

	hasher := internal.NewHasher(lastReader)

	writer, err := internal.NewFileWriter(hasher.Reader(), outPath, fileReader.Type())
	if err != nil {
		return err
	}

	err = writer.Run(appCtx.noCopy)
	hasher.After(asset)

	return nil
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

func noExt(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}
