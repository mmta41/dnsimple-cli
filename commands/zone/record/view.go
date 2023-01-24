package record

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/mmta41/dnsimple-cli/commands/cmdutils"
	"github.com/spf13/cobra"
	"strconv"
)

func NewCmdViewRecorc(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "record <zoneName> <zoneId>",
		Args:    cmdutils.MinimumArgs(2, "no zone name or record id specified"),
		Short:   "View specific zone's record.",
		Long:    heredoc.Doc(`Show specific record belong to specific zone`),
		PreRunE: f.PreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			parseInt, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}
			return ViewRecord(f, args[0], parseInt)
		},
	}

	return cmd
}

func ViewRecord(f *cmdutils.Factory, zoneName string, zoneId int64) error {
	fmt.Fprintf(f.IO.StdOut, "Record[%v] of %v\n", zoneId, zoneName)
	fmt.Fprintf(f.IO.StdOut, "This function is not implemented\n")
	return nil
}
