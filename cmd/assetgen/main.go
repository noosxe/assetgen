package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	code := run(os.Stdout, os.Stderr, os.Args)
	if code != 0 {
		os.Exit(code)
	}
}

const usage = `usage: assetgen <command>

assetgen - generate manifest files for your static assets

commands:
	generate	Generates a manifest file
	version		Prints the version
`

func run(stdout, stderr io.Writer, args []string) (code int) {
	if len(args) < 2 {
		fmt.Fprint(stderr, usage)
		return 64
	}
	switch args[1] {
	case "generate":
		return generate(args[2:])
	case "help", "-help", "--help", "-h":
		fmt.Fprint(stdout, usage)
		return 0
	}
	fmt.Fprint(stderr, usage)
	return 64
}

const usageGenerate = `usage: assetgen generate [<args>...]

Generates the manifest file

Args:
	-config
		Config file path. (default ./assetgen.yaml)
	-out
		Specify the output directory.
`

func generate(args []string) (code int) {
	cmd := flag.NewFlagSet("generate", flag.ExitOnError)

	configFlag := cmd.String("config", "./assetgen.yaml", "config file path")
	outFlag := cmd.String("out", "", "output directory")

	cmd.Usage = func() { fmt.Print(usageGenerate) }
	err := cmd.Parse(args)
	if err != nil {
		cmd.PrintDefaults()
		return 64
	}
	ctx := AppContext{configPath: *configFlag, outPath: *outFlag}
	err = normalizeContext(&ctx)
	if err != nil {
		return 1
	}

	return GenerateManifest(ctx)
}

func normalizeContext(ctx *AppContext) error {
	configPath, err := filepath.Abs(ctx.configPath)
	if err != nil {
		return err
	}
	ctx.configPath = configPath
	ctx.configDir = filepath.Dir(ctx.configPath)

	if ctx.outPath != "" {
		if !filepath.IsAbs(ctx.outPath) {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			ctx.outPath = filepath.Join(wd, ctx.outPath)
		}
	}

	return nil
}
