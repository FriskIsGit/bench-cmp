package main

import (
	"fmt"
	"os"
)

const Version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		PrintHelp(Version)
		return
	}

	args := ParseArguments(os.Args[1:], Version)
	if !args.Validate() {
		return
	}

	if args.ListTests {
		PrintTests()
		return
	}

	if err := execute(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
