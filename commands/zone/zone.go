package zone

import (
	"github.com/mmta41/dnsimple-cli/commands/cmdutils"
	recordCmd "github.com/mmta41/dnsimple-cli/commands/zone/record"
	"github.com/spf13/cobra"
)

func NewCmdZone(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone <command>",
		Short: "Manage dns zone",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdView(f))
	cmd.AddCommand(recordCmd.NewCmdViewRecorc(f))

	return cmd
}
