package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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

const usageGenerate = `usage: assetgen generate [<args>...]

Generates the manifest file

Args:
	-config
		Config file path. (default ./assetgen.yaml)
	-out
		Specify the output directory. (default ./dist)
`

func generate(args []string) (code int) {
	cmd := flag.NewFlagSet("generate", flag.ExitOnError)

	var config string
	cmd.StringVar(&config, "config", "./assetgen.yaml", "config file path")

	var out string
	cmd.StringVar(&out, "out", "./dist", "output directory")

	cmd.Usage = func() { fmt.Print(usageGenerate) }
	err := cmd.Parse(args)
	if err != nil {
		cmd.PrintDefaults()
		return 64
	}

	return GenerateManifest(config, out)
}
