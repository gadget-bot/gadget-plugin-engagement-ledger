package plugin

import (
	"github.com/gadget-bot/gadget-plugin-engagement-ledger/internal/slackclient"
	"github.com/gadget-bot/gadget/core"
	"github.com/gadget-bot/gadget/router"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// SchedulerRegistrar is implemented by gadget-plugin-scheduler (or any compatible scheduler).
type SchedulerRegistrar interface {
	RegisterJob(name, cronExpr string, fn func()) error
}

// ChatEventSource is implemented by gadget-plugin-chat-adapters (or any compatible adapter).
type ChatEventSource interface {
	OnMessage(fn func(workspaceID, userID, channelID, ts string))
}

// Plugin holds the assembled plugin state.
type Plugin struct {
	cfg       Config
	scheduler SchedulerRegistrar // nil when integration disabled
	chatSrc   ChatEventSource    // nil when integration disabled
}

// Option is a functional option for Plugin.
type Option func(*Plugin)

// WithScheduler wires a SchedulerRegistrar for monthly awards and leaderboard jobs.
func WithScheduler(s SchedulerRegistrar) Option { return func(p *Plugin) { p.scheduler = s } }

// WithChatAdapters wires a ChatEventSource for normalized message events.
func WithChatAdapters(src ChatEventSource) Option { return func(p *Plugin) { p.chatSrc = src } }

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
	bot.Router.AddSlashCommandRoutes(p.slashCommandRoutes(bot))

	if p.scheduler != nil {
		p.registerScheduledJobs()
	}

	if p.chatSrc != nil {
		p.subscribeChatEvents()
	}
}

// channelMessageRoutes returns routes for channel message events.
func (p *Plugin) channelMessageRoutes(bot *core.Gadget) []router.ChannelMessageRoute {
	client := slackclient.Wrap(bot.Client)
	_ = client // available for handler injection in implementing issues

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
				// TODO(issue #5): enforce eligibility rules
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
				// TODO(issue #6): add playful Penny responses for --
				log.Debug().Str("route", "engagement.quip.decrement").Msg("stub: -- quip handler")
			},
		},
	}
}

// slashCommandRoutes returns routes for slash command events.
func (p *Plugin) slashCommandRoutes(bot *core.Gadget) []router.SlashCommandRoute {
	client := slackclient.Wrap(bot.Client)
	_ = client // available for handler injection in implementing issues

	command := "/give"
	if p.cfg.Commands.Prefix != "" {
		command = "/" + p.cfg.Commands.Prefix + "-give"
	}

	return []router.SlashCommandRoute{
		{
			Route: router.Route{
				Name:        "engagement.award.give",
				Description: "Award a point to a user via slash command",
			},
			Command: command,
			Plugin: func(ctx router.HandlerContext, cmd slack.SlashCommand) {
				// TODO(issue #4): implement command routing with configurable top_level/subcommand modes
				log.Debug().Str("route", "engagement.award.give").Msg("stub: /give handler")
			},
		},
	}
}

// registerScheduledJobs registers monthly award and leaderboard jobs with the scheduler.
func (p *Plugin) registerScheduledJobs() {
	if p.cfg.MonthlyAward.Enabled {
		if err := p.scheduler.RegisterJob("engagement.monthly_award", "@monthly", func() {
			// TODO(issue #8): implement monthly award job with timezone boundary logic
			log.Debug().Msg("stub: monthly award job")
		}); err != nil {
			log.Error().Err(err).Msg("Failed to register monthly award job")
		}
	}

	if p.cfg.Leaderboard.Enabled {
		schedule := p.cfg.Leaderboard.Schedule
		if err := p.scheduler.RegisterJob("engagement.leaderboard", schedule, func() {
			// TODO(issue #9): implement weekly leaderboard publishing
			log.Debug().Msg("stub: leaderboard job")
		}); err != nil {
			log.Error().Err(err).Msg("Failed to register leaderboard job")
		}
	}
}

// subscribeChatEvents subscribes to normalized message events from chat-adapters.
func (p *Plugin) subscribeChatEvents() {
	p.chatSrc.OnMessage(func(workspaceID, userID, channelID, ts string) {
		// TODO(issue #7): track active-user activity via message event listener
		log.Debug().Str("workspace", workspaceID).Str("user", userID).Msg("stub: chat adapter message event")
	})
}
