package zone

import (
	"context"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/mmta41/dnsimple-cli/commands/cmdutils"
	"github.com/mmta41/dnsimple-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

func NewCmdList(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Args:  cobra.ExactArgs(0),
		Short: "Lists the zones in the account.",
		Long: heredoc.Doc(`A DNS zone is an administrative space which allows for more granular control of DNS components,
such as authoritative nameservers. The domain name space is a hierarchical tree,with the DNS root domain at the top.
A DNS zone starts at a domain within the tree and can also extend down into subdomains so that multiple subdomains can be managed by one entity.`),
		PreRunE: f.PreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listZone(f)
		},
	}

	return cmd
}

func listZone(f *cmdutils.Factory) error {
	client, err := f.Client()
	if err != nil {
		return err
	}

	zones, err := client.Zones.ListZones(context.Background(), client.AccountIdStr(), &dnsimple.ZoneListOptions{})
	if err != nil {
		return client.ConvertError(err)
	}

	var options = make([]string, len(zones.Data))

	for i, zone := range zones.Data {
		options[i] = fmt.Sprintf("[%v] %v (reverse: %v)", zone.ID, zone.Name, zone.Reverse)
	}
	if f.IO.PromptEnabled() {
		var selected int
		err := prompt.SelectOne(&selected, "Select DNS Zone:", options)
		if err != nil {
			return err
		}
		if selected >= 0 {
			return viewZone(f, zones.Data[selected].Name)
		}
		return nil
	}
	for _, opt := range options {
		fmt.Fprintf(f.IO.StdOut, "%v\n", opt)
	}
	return nil
}
