package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

func main() {

	// Only take in the source and destination directories
	parser := argparse.NewParser("dirmap", "Maps a directory onto a new one given a set of transformations.")
	src_dir := parser.StringPositional(&argparse.Options{Required: false, Help: "Source directory to transform"})
	dst_dir := parser.StringPositional(&argparse.Options{Required: false, Help: "Destination directory to write to"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println("Couldn't parse arguments!")
		os.Exit(-1)
	}

	// Some sainity checks
	if src_dir != nil {

		stat, err := os.Stat(*src_dir)
		if err != nil {
			fmt.Println("Cannot stat directory", *src_dir)
			os.Exit(0)
		} else if !stat.IsDir() {
			fmt.Println("Source directory is not actually a directory!")
			os.Exit(0)
		}

	} else { // Otherwise we default to the current working directory

		tmp, err := os.Getwd()
		if err != nil {
			fmt.Println("Cannot get current working directory :(")
			os.Exit(00)
		} else {
			src_dir = &tmp
		}

	}

	if dst_dir != nil {

		stat, err := os.Stat(*dst_dir)
		if err != nil {
			fmt.Println("Cannot stat directory", *dst_dir)
			os.Exit(0)
		} else if !stat.IsDir() {
			fmt.Println("Destination directory is not actually a directory!")
			os.Exit(0)
		}

	} else { // If we didn't get a destination directory, then we default to the current working directory

		tmp, err := os.Getwd()
		if err != nil {
			fmt.Println("Cannot get current working directory :(")
			os.Exit(0)
		}

		dst_dir = &tmp

	}

	// Now that the preamble is done, we can work on loading the TUI
	err = StartTui(*src_dir, *dst_dir)
	if err != nil {
		fmt.Println("Got error:", err.Error())
		os.Exit(-1)
	}

	os.Exit(0)

}
