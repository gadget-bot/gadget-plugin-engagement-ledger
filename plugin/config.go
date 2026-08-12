package plugin

import "os"

// Config holds all v1 plugin configuration.
type Config struct {
	FeedbackReaction string // emoji reaction posted on successful award; default: "+1"
}

// ConfigFromEnv returns a Config populated from environment variables, with defaults applied.
func ConfigFromEnv() Config {
	return Config{
		FeedbackReaction: getEnvOrDefault("ENGAGEMENT_FEEDBACK_REACTION", "+1"),
	}
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
