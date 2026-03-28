package plugin

import (
	"os"
	"strconv"
)

// CommandsMode controls how the plugin registers its commands.
type CommandsMode string

const (
	CommandsModeTopLevel   CommandsMode = "top_level"
	CommandsModeSubcommand CommandsMode = "subcommand"
)

// TimezoneMode controls the timezone boundary for monthly award calculations.
type TimezoneMode string

const (
	TimezoneModeWorkspaceLocal TimezoneMode = "workspace_local"
	TimezoneModeUTC            TimezoneMode = "utc"
)

// Config holds all v1 plugin configuration.
type Config struct {
	Commands struct {
		Mode   CommandsMode // default: top_level
		Prefix string
	}
	Feedback struct {
		Reaction string // default: "+1"
	}
	MonthlyAward struct {
		Enabled      bool
		Points       int          // default: 10
		TimezoneMode TimezoneMode // default: utc
	}
	Leaderboard struct {
		Enabled   bool
		ChannelID string
		Schedule  string // default: "0 8 * * 1"
	}
	Integrations struct {
		Scheduler    struct{ Enabled bool }
		ChatAdapters struct{ Enabled bool }
	}
}

// ConfigFromEnv returns a Config populated from environment variables, with defaults applied.
func ConfigFromEnv() Config {
	var cfg Config

	cfg.Commands.Mode = CommandsMode(getEnvOrDefault("ENGAGEMENT_COMMANDS_MODE", string(CommandsModeTopLevel)))
	cfg.Commands.Prefix = os.Getenv("ENGAGEMENT_COMMANDS_PREFIX")

	cfg.Feedback.Reaction = getEnvOrDefault("ENGAGEMENT_FEEDBACK_REACTION", "+1")

	cfg.MonthlyAward.Enabled = parseBoolEnv("ENGAGEMENT_MONTHLY_AWARD_ENABLED")
	cfg.MonthlyAward.Points = parseIntEnvOrDefault("ENGAGEMENT_MONTHLY_AWARD_POINTS", 10)
	cfg.MonthlyAward.TimezoneMode = TimezoneMode(getEnvOrDefault("ENGAGEMENT_MONTHLY_AWARD_TIMEZONE_MODE", string(TimezoneModeUTC)))

	cfg.Leaderboard.Enabled = parseBoolEnv("ENGAGEMENT_LEADERBOARD_ENABLED")
	cfg.Leaderboard.ChannelID = os.Getenv("ENGAGEMENT_LEADERBOARD_CHANNEL_ID")
	cfg.Leaderboard.Schedule = getEnvOrDefault("ENGAGEMENT_LEADERBOARD_SCHEDULE", "0 8 * * 1")

	cfg.Integrations.Scheduler.Enabled = parseBoolEnv("ENGAGEMENT_INTEGRATIONS_SCHEDULER_ENABLED")
	cfg.Integrations.ChatAdapters.Enabled = parseBoolEnv("ENGAGEMENT_INTEGRATIONS_CHAT_ADAPTERS_ENABLED")

	return cfg
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func parseBoolEnv(key string) bool {
	v := os.Getenv(key)
	if v == "" {
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}

func parseIntEnvOrDefault(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}
