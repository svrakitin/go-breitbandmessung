package breitbandmessung

import (
	"context"

	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "breitbandmessung",
	}
	cmd.AddCommand(newSnapshotCommand())

	return cmd
}

func Execute(ctx context.Context) error {
	cmd := newRootCommand()
	return cmd.ExecuteContext(ctx)
}
