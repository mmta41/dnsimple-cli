package version

import (
	"fmt"
	"github.com/mmta41/dnsimple-cli/pkg/iostreams"
	"strings"

	"github.com/spf13/cobra"
)

func NewCmdVersion(s *iostreams.IOStreams, version, buildDate string) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:     "version",
		Short:   "show dns-cli version information",
		Long:    ``,
		Aliases: []string{"v"},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(s.StdOut, Scheme(version, buildDate))
			return nil
		},
	}
	return versionCmd
}

func Scheme(version, buildDate string) string {
	version = strings.TrimPrefix(version, "v")

	if buildDate != "" {
		version = fmt.Sprintf("%s (%s)", version, buildDate)
	}

	return fmt.Sprintf("dnsimple-cli version %s\n", version)
}
