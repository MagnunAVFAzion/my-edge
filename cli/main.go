package main

import (
	"cli/build"
	"cli/pack"
	"cli/send"
	"fmt"
	"os"
)

func clean() error {
	err := os.RemoveAll("./out")
	if err != nil {
		panic(fmt.Errorf("error in clean: %w", err))
	}

	return nil
}

func main() {
	clean()

	err := build.RunVercelBuild()
	if err != nil {
		panic(fmt.Errorf("error running vercel build: %w", err))
	}

	err = pack.PackProject()
	if err != nil {
		panic(fmt.Errorf("error packing project: %w", err))
	}

	send.SendProject()

	fmt.Println("Done!")
}
