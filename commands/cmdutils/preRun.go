package cmdutils

import "github.com/spf13/cobra"

func (f *Factory) PreRunE(cmd *cobra.Command, args []string) error {
	if !f.NeedPrompt() {
		return nil
	}
	return f.PromptConfig()
}
