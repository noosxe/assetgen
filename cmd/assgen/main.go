package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func main() {
	code := run(os.Stdout, os.Stderr, os.Args)
	if code != 0 {
		os.Exit(code)
	}
}

const usage = `usage: assgen <command>

assgen - generate manifest files for your static assets

commands:
	generate	Generates a manifest file
	version		Prints the version
`

func run(stdout, stderr io.Writer, args []string) (code int) {
	if len(args) < 2 {
		fmt.Fprint(stderr, usage)
		return 64 // EX_USAGE
	}
	switch args[1] {
	case "generate":
		return generate(args[2:])
	case "help", "-help", "--help", "-h":
		fmt.Fprint(stdout, usage)
		return 0
	}
	fmt.Fprint(stderr, usage)
	return 64 // EX_USAGE
}

type Config struct {
	Styles  []string `yaml:"styles"`
	Scripts []string `yaml:"scripts"`
	Out     *string  `yaml:"out"`
}

const usageGenerate = `usage: assgen generate [<args>...]

Generates the manifest file

Args:
	-config
		Config file path. (default ./assgen.yaml)
	-out
		Specify the output directory. (default ./dist)
`

func generate(args []string) (code int) {
	cmd := flag.NewFlagSet("generate", flag.ExitOnError)

	var config string
	cmd.StringVar(&config, "config", "./assgen.yaml", "config file path")

	var out string
	cmd.StringVar(&out, "out", "./dist", "output directory")

	cmd.Usage = func() { fmt.Print(usageGenerate) }
	err := cmd.Parse(args)
	if err != nil {
		cmd.PrintDefaults()
		return 64
	}

	configFilePath, err := filepath.Abs(config)
	if err != nil {
		log.Fatalln(err)
	}

	configFileDir := filepath.Dir(configFilePath)
	outputPath := filepath.Join(configFileDir, out)

	log.Println("ensuring output path")
	_, err = os.Stat(outputPath)
	if os.IsNotExist(err) {
		log.Println("output path does not exist, creating...")
		err := os.Mkdir(outputPath, 0755)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("output path created")
	} else {
		log.Println("output path exists")
	}

	configDoc, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
		return 1
	}

	c := Config{}

	err = yaml.Unmarshal(configDoc, &c)
	if err != nil {
		log.Fatal(err)
	}

	scripts := c.Scripts
	styles := c.Styles

	log.Println("processing scripts")
	err = processGlobs(scripts, configFileDir, outputPath)

	log.Println("processing styles")
	err = processGlobs(styles, configFileDir, outputPath)

	return 0
}

func processGlobs(globs []string, configFileDir string, outputPath string) error {
	for _, script := range globs {
		fullpath := filepath.Join(configFileDir, script)
		matches, err := filepath.Glob(fullpath)
		if err != nil {
			return err
		}

		for _, match := range matches {
			rel, err := filepath.Rel(configFileDir, match)
			if err != nil {
				return err
			}

			destPath := filepath.Join(outputPath, rel)
			log.Printf("copying %s\n", rel)
			hash, err := copyFile(match, destPath)
			if err != nil {
				return err
			}

			fmt.Println(hash)
		}
	}

	return nil
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
