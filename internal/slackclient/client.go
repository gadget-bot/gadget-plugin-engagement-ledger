package slackclient

import "github.com/slack-go/slack"

// Client is a minimal interface over *slack.Client covering all v1 Slack calls.
// Handlers accept Client rather than *slack.Client to enable mock injection in tests.
type Client interface {
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
	AddReaction(name string, item slack.ItemRef) error
	GetUserInfo(user string) (*slack.User, error)
}

// Wrap adapts *slack.Client to Client. *slack.Client satisfies the interface
// structurally; Wrap is a named call-site for explicitness.
func Wrap(c *slack.Client) Client { return c }
