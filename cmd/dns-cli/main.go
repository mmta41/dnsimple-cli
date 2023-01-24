package main

import (
	"errors"
	"fmt"
	surveyCore "github.com/AlecAivazis/survey/v2/core"
	"github.com/mgutz/ansi"
	"github.com/mmta41/dnsimple-cli/commands"
	"github.com/mmta41/dnsimple-cli/commands/cmdutils"
	"github.com/mmta41/dnsimple-cli/commands/help"
	"github.com/mmta41/dnsimple-cli/pkg/iostreams"
	"github.com/mmta41/dnsimple-cli/pkg/tableprinter"
	"github.com/spf13/cobra"
	"net"
	"os"
	"strings"
)

var version = "DEV"
var build string

const debug = true

func main() {

	cmdFactory := cmdutils.NewFactory(debug)

	if !cmdFactory.IO.ColorEnabled() {
		surveyCore.DisableColor = true
	} else {
		// Override survey's choice of color for default values
		// For default values for e.g. `Input` prompts, Survey uses the literal "white" color,
		// which makes no sense on dark terminals and is literally invisible on light backgrounds.
		// This overrides Survey to output a gray color for 256-color terminals and "default" for basic terminals.
		surveyCore.TemplateFuncsWithColor["color"] = func(style string) string {
			switch style {
			case "white":
				if cmdFactory.IO.Is256ColorSupported() {
					return fmt.Sprintf("\x1b[%d;5;%dm", 38, 242)
				}
				return ansi.ColorCode("default")
			default:
				return ansi.ColorCode(style)
			}
		}
	}

	cmdList := commands.InitCommands(cmdFactory, version, build)

	cfg, err := cmdFactory.Config()
	if err != nil {
		cmdFactory.IO.Logf("failed to read configuration:  %s\n", err)
		os.Exit(2)
	}

	if cfg.NoPrompt != "" {
		cmdFactory.IO.SetPrompt(cfg.NoPrompt)
	}

	if forceHyperlinks := os.Getenv("FORCE_HYPERLINKS"); forceHyperlinks != "" && forceHyperlinks != "0" {
		cmdFactory.IO.SetDisplayHyperlinks("always")
	} else if cfg.ForceHyperLinks == "true" {
		cmdFactory.IO.SetDisplayHyperlinks("auto")
	}

	var expandedArgs []string
	if len(os.Args) > 0 {
		expandedArgs = os.Args[1:]
	}

	// Override the default column separator of tableprinter to double spaces
	tableprinter.SetTTYSeparator("  ")
	// Override the default terminal width of tableprinter
	tableprinter.SetTerminalWidth(cmdFactory.IO.TerminalWidth())
	// set whether terminal is a TTY or non-TTY
	tableprinter.SetIsTTY(cmdFactory.IO.IsOutputTTY())

	cmdList.SetArgs(expandedArgs)

	if cmd, err := cmdList.ExecuteC(); err != nil {
		printError(cmdFactory.IO, err, cmd, debug, true)
	}

	if help.HasFailed() {
		os.Exit(1)
	}
}

func printError(streams *iostreams.IOStreams, err error, cmd *cobra.Command, debug, shouldExit bool) {
	if errors.Is(err, cmdutils.SilentError) {
		return
	}
	color := streams.Color()
	printMore := true
	exitCode := 1

	var dnsError *net.DNSError
	if errors.As(err, &dnsError) {
		streams.Logf("%s error connecting to %s\n", color.FailedIcon(), dnsError.Name)
		if debug {
			streams.Log(color.FailedIcon(), dnsError)
		}
		streams.Logf("%s check your internet connection or status of dnsimple.com\n", color.DotWarnIcon())
		printMore = false
	}
	if printMore {
		var exitError *cmdutils.ExitError
		if errors.As(err, &exitError) {
			streams.Logf("%s %s %s=%s\n", color.FailedIcon(), color.Bold(exitError.Details), color.Red("error"), exitError.Err)
			exitCode = exitError.Code
			printMore = false
		}

		if printMore {
			streams.Log(err)

			var flagError *cmdutils.FlagError
			if errors.As(err, &flagError) || strings.HasPrefix(err.Error(), "unknown command ") {
				if !strings.HasSuffix(err.Error(), "\n") {
					streams.Log()
				}
				streams.Log(cmd.UsageString())
			}
		}
	}

	if cmd != nil {
		cmd.Print("\n")
	}
	if shouldExit {
		os.Exit(exitCode)
	}
}
