package plugin

import (
	"github.com/gadget-bot/gadget-plugin-engagement-ledger/internal/slackclient"
	"github.com/gadget-bot/gadget/core"
	"github.com/gadget-bot/gadget/router"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack/slackevents"
)

// Plugin holds the assembled plugin state.
type Plugin struct {
	cfg Config
}

// Option is a functional option for Plugin. No options exist in v1; reserved for v2/v3 integrations.
type Option func(*Plugin)

// Register creates a Plugin, applies options, and registers all routes and jobs on bot.
func Register(bot *core.Gadget, cfg Config, opts ...Option) {
	p := &Plugin{cfg: cfg}
	for _, o := range opts {
		o(p)
	}
	p.register(bot)
}

func (p *Plugin) register(bot *core.Gadget) {
	bot.Router.AddChannelMessageRoutes(p.channelMessageRoutes(bot))
}

// channelMessageRoutes returns routes for channel message events.
func (p *Plugin) channelMessageRoutes(bot *core.Gadget) []router.ChannelMessageRoute {
	client := slackclient.Wrap(bot.Client)
	_ = client // pass to handler closures when implementing

	return []router.ChannelMessageRoute{
		{
			Route: router.Route{
				Name:        "engagement.award.mention",
				Pattern:     `<@[A-Z0-9]+>\s*\+\+`,
				Description: "Award a point to a mentioned user via @user++ syntax",
				Priority:    10,
			},
			Plugin: func(ctx router.HandlerContext, ev slackevents.MessageEvent, message string) {
				// TODO(issue #3): implement mention parser for @user++ with recipient dedupe
				// TODO(issue #4): enforce eligibility rules
				log.Debug().Str("route", "engagement.award.mention").Msg("stub: @user++ handler")
			},
		},
		{
			Route: router.Route{
				Name:        "engagement.quip.decrement",
				Pattern:     `<@[A-Z0-9]+>\s*--`,
				Description: "Playful response when -- is used (points are never removed)",
				Priority:    10,
			},
			Plugin: func(ctx router.HandlerContext, ev slackevents.MessageEvent, message string) {
				// TODO(issue #5): add playful Penny responses for --
				log.Debug().Str("route", "engagement.quip.decrement").Msg("stub: -- quip handler")
			},
		},
	}
}
