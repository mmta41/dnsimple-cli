package config

import (
	"github.com/mmta41/dnsimple-cli/pkg/prompt"
	"os"
)

type Account struct {
	Token string `yaml:"token"`
	Id    int64  `yaml:"id"`
}
type Config struct {
	Account
	NoPrompt        string `yaml:"no_prompt"`
	ForceHyperLinks string `yaml:"force_hyperlinks"`
}

func (cfg *Config) Save() error {
	return WriteConfigFile(cfg)
}

func Init(opt *Account) (*Config, error) {
	cfg, err := readFromFile()
	if cfg == nil {
		cfg = &Config{}
	}

	if err != nil && !os.IsNotExist(err) {
		return cfg, err
	}
	overrideEnv(cfg)

	//override opt-in
	if opt.Token != "" {
		cfg.Token = opt.Token
	}

	if opt.Id != 0 {
		cfg.Id = opt.Id
	}

	return cfg, nil
}

func Prompt(question, defaultVal string) (envVal string, err error) {
	err = prompt.AskQuestionWithInput(&envVal, "config", question, defaultVal, false)
	if err != nil {
		return
	}
	return
}

func PromptMulti(question string, options []string) (envVal string, err error) {
	err = prompt.Select(&envVal, "config", question, options)
	if err != nil {
		return
	}
	return
}

func ConfirmSave(cfg *Config) error {
	var result bool
	err := prompt.Confirm(&result, "Do you want to store this credentials?", false)
	if err != nil {
		return err
	}
	if !result {
		return nil
	}

	err = WriteConfigFile(cfg)
	if err != nil {
		return err
	}

	return nil
}
