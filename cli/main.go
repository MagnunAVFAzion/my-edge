package main

import (
	"cli/build"
	"cli/pack"
	"cli/send"
	"fmt"
	"os"
)

func clean() error {
	err := os.RemoveAll("out")
	if err != nil {
		panic(fmt.Errorf("error in clean: %w", err))
	}

	err = os.Remove("entry.js")
	if err != nil {
		panic(fmt.Errorf("error in clean: %w", err))
	}

	return nil
}

func main() {
	clean()

	err := build.Exec()
	if err != nil {
		panic(fmt.Errorf("error running build: %w", err))
	}

	err = pack.PackProject()
	if err != nil {
		panic(fmt.Errorf("error packing project: %w", err))
	}

	err = send.SendProject()
	if err != nil {
		panic(fmt.Errorf("error sending project: %w", err))
	}

	fmt.Println("\nSuccess!")
}
