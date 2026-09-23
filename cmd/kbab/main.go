package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "kbab",
		Short: "Generate production-grade Dockerfiles and .dockerignore files",
	}

	create := &cobra.Command{
		Use:   "create [path]",
		Short: "Detect the project and generate container setup",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			fmt.Fprintln(cmd.OutOrStdout(), "create: ", path)
			return nil
		},
	}

	root.AddCommand(create)
	return root
}
