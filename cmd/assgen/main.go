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
		return generate(stdout, stderr, args[2:])
	case "help", "-help", "--help", "-h":
		fmt.Fprint(stdout, usage)
		return 0
	}
	fmt.Fprint(stderr, usage)
	return 64 // EX_USAGE
}

const usageGenerate = `usage: assgen generate [<args>...]

Args:
	-path <path>
		Generates manifest for all assets in path.
`

func generate(stdout, stderr io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("generate", flag.ExitOnError)

	pathFlag := cmd.String("path", "", "")

	err := cmd.Parse(args)
	if err != nil {
		fmt.Fprint(stderr, usageGenerate)
	}

	fmt.Println(*pathFlag)

	return 0
}
