package config

import (
	"os"
	"strconv"
)

func overrideEnv(cfg *Config) {
	token := os.Getenv("DNS_TOKEN")
	accountId := os.Getenv("DNS_ACCOUNT")

	if token != "" {
		cfg.Token = token
	}

	if accountId != "" {
		parseInt, err := strconv.ParseInt(accountId, 10, 64)
		if err != nil {
			return
		}
		cfg.Id = parseInt
	}
}
