package main

import (
	"fmt"
	"os"

	"github.com/santiagozsh/kbab/internal/detector"
)

func main() {
	targetDir := "."
	if len(os.Args) > 1 {
		targetDir = os.Args[1]
	}

	projectType, err := detector.Detect(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("project detector success!!: %s\n", projectType)
}
