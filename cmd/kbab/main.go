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

	cfg, err := detector.Detect(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🚀 Project detected successfully!")
	fmt.Printf("  Type:            %s\n", cfg.Type)
	fmt.Printf("  Runtime Version: %s\n", cfg.RuntimeVersion)
	fmt.Printf("  Package Manager: %s\n", cfg.PackageManager)
}
