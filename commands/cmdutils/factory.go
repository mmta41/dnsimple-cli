package cmdutils

import (
	"github.com/mmta41/dnsimple-cli/internal/api"
	"github.com/mmta41/dnsimple-cli/internal/config"
	"github.com/mmta41/dnsimple-cli/pkg/iostreams"
	"strconv"
)

type Factory struct {
	config      *config.Config
	configError error
	IO          *iostreams.IOStreams
	client      *api.Client
	OptIn       *config.Account
	Debug       bool
}

func (f *Factory) Config() (*config.Config, error) {
	if f.config.Token != "" {
		return f.config, f.configError
	}
	f.config, f.configError = config.Init(f.OptIn)
	return f.config, f.configError
}

func (f *Factory) PromptConfig() error {
	var err error
	f.config = &config.Config{}
	f.config.Token, err = config.Prompt("Please enter your dnsimple token: ", f.config.Token)
	if err != nil {
		return err
	}
	client, err := f.Client()
	if err != nil {
		return err
	}
	accounts, err := client.Login()
	if err != nil {
		return err
	}

	options := make([]string, len(accounts))

	for i, acc := range accounts {
		options[i] = "[" + strconv.FormatInt(acc.ID, 10) + "] " + acc.Email + " (" + acc.PlanIdentifier + ")"
	}

	selected, err2 := config.PromptMulti("Selected Account:\n", options)
	if err2 != nil {
		return err2
	}

	for i, acc := range options {
		if acc == selected {
			f.config.Id = accounts[i].ID
			break
		}
	}

	err = config.ConfirmSave(f.config)
	if err != nil {
		return err
	}
	return nil
}

func (f *Factory) NeedPrompt() bool {
	cfg, err := f.Config()
	if err != nil || cfg == nil {
		return true
	}
	return cfg.Token == ""
}

func (f *Factory) Client() (*api.Client, error) {
	if f.client != nil {
		return f.client, nil
	}
	cfg, err := f.Config()
	if err != nil {
		return nil, err
	}
	f.client = api.InitClient(cfg.Token, f.Debug)
	f.client.SetAccountId(f.config.Id)
	return f.client, nil
}

func NewFactory(debug bool) *Factory {
	return &Factory{
		OptIn:  &config.Account{},
		config: &config.Config{},
		IO:     iostreams.Init(),
		Debug:  debug,
	}
}
