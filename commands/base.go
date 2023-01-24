package commands

import (
	"errors"
	"github.com/MakeNowJust/heredoc"
	"github.com/mmta41/dnsimple-cli/commands/cmdutils"
	completionCmd "github.com/mmta41/dnsimple-cli/commands/completion"
	"github.com/mmta41/dnsimple-cli/commands/help"
	versionCmd "github.com/mmta41/dnsimple-cli/commands/version"
	zoneCmd "github.com/mmta41/dnsimple-cli/commands/zone"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func InitCommands(f *cmdutils.Factory, version, buildDate string) *cobra.Command {
	c := f.IO.Color()

	var cmd = &cobra.Command{
		Use:           "dns-cli <command> <subcommand> [flags]",
		Short:         "A dnsimple.com CLI/UI Tool",
		Long:          `dns-cli is an open source dnsimple.com CLI tool bringing DNSimple API to your command line`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Annotations: map[string]string{
			"help:environment": heredoc.Doc(`
			DNS_TOKEN: an authentication token for API requests. Setting this avoids being
			prompted to authenticate and overrides any previously stored credentials.
			Can be set in the config with 'dns-cli config set token xxxxxx'
			
			DNS_ACCOUNT: account id for API request. setting this voids being prompted.
			Can be set in the config with 'dns-cli config set account xx'

			NO_PROMPT: set to 1 (true) or 0 (false) to disable and enable prompts respectively

			NO_COLOR: set to any value to avoid printing ANSI escape sequences for color output.

			FORCE_HYPERLINKS: set to 1 to force hyperlinks to be output, even when not outputting to a TTY

			DNS_CONFIG_DIR: set to a directory path to override the global configuration location 
		`),
			"help:feedback": heredoc.Docf(`
			Encountered a bug or want to suggest a feature?
			contact us: %s`, c.Bold(c.Yellow("info@dnsimple.com"))),
		},
	}

	cmd.SetOut(f.IO.StdOut)
	cmd.SetErr(f.IO.StdErr)

	cmd.PersistentFlags().Bool("help", false, "Show help for command")
	cmd.SetHelpFunc(func(command *cobra.Command, args []string) {
		help.RootHelpFunc(f.IO.Color(), command, args)
	})
	cmd.SetUsageFunc(help.RootUsageFunc)
	cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		if errors.Is(err, pflag.ErrHelp) {
			return err
		}
		return &cmdutils.FlagError{Err: err}
	})

	formattedVersion := versionCmd.Scheme(version, buildDate)
	cmd.SetVersionTemplate(formattedVersion)
	cmd.Version = formattedVersion

	cmd.AddCommand(completionCmd.NewCmdCompletion(f.IO))

	cmd.AddCommand(zoneCmd.NewCmdZone(f))

	cmd.Flags().StringVarP(&f.OptIn.Token, "token", "t", "", "Your dnsimple access token")
	cmd.Flags().Int64VarP(&f.OptIn.Id, "id", "i", 0, "Your dnsimple account id")
	cmd.Flags().BoolP("version", "v", false, "show dns-cli version information")
	return cmd
}
