package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

func Execute() error {
	return newRootCmd(os.Stdout, os.Stderr).Execute()
}

func newRootCmd(stdout io.Writer, stderr io.Writer) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "lexi",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	return rootCmd
}
