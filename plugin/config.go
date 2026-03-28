package plugin

import "os"

// Config holds all v1 plugin configuration.
type Config struct {
	Feedback struct {
		Reaction string // default: "+1"
	}
}

// ConfigFromEnv returns a Config populated from environment variables, with defaults applied.
func ConfigFromEnv() Config {
	var cfg Config

	cfg.Feedback.Reaction = getEnvOrDefault("ENGAGEMENT_FEEDBACK_REACTION", "+1")

	return cfg
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
