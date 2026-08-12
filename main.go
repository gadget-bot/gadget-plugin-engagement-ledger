package main

import (
	"github.com/gadget-bot/gadget-plugin-engagement-ledger/plugin"
	"github.com/gadget-bot/gadget/core"
	"github.com/rs/zerolog/log"
)

func main() {
	bot, err := core.Setup()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to set up bot")
	}

	plugin.Register(bot, plugin.ConfigFromEnv())
	if err := bot.Run(); err != nil {
		log.Fatal().Err(err).Msg("bot stopped")
	}
}
