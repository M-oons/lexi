package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/m-oons/lexi/cmd/random"
	"github.com/spf13/cobra"
)

func Execute() error {
	root := newRootCmd(os.Stdout, os.Stderr)
	return root.Execute()
}

func newRootCmd(stdout io.Writer, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "lexi",
		Args:          cobra.MinimumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("unknown command %q", args[0])
		},
	}
	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	cmd.AddCommand(random.NewCmd())

	return cmd
}
