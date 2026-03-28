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

	cfg := plugin.ConfigFromEnv()
	var opts []plugin.Option

	if cfg.Integrations.Scheduler.Enabled {
		// opts = append(opts, plugin.WithScheduler(myScheduler))
		log.Warn().Msg("scheduler integration enabled but not wired — skipping")
	}
	if cfg.Integrations.ChatAdapters.Enabled {
		// opts = append(opts, plugin.WithChatAdapters(myAdapter))
		log.Warn().Msg("chat_adapters integration enabled but not wired — skipping")
	}

	plugin.Register(bot, cfg, opts...)
	if err := bot.Run(); err != nil {
		log.Fatal().Err(err).Msg("bot stopped")
	}
}
