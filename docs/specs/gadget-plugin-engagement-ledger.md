# gadget-plugin-engagement-ledger Specification

## Purpose
Provide a standalone engagement points system with an immutable ledger and user-to-user point transfers via mention syntax.

## Standalone and Optional Integrations
- Operates standalone with direct command/API input.
- Optionally registers recurring jobs with `gadget-plugin-scheduler` (required for monthly awards and weekly leaderboard). _(v3)_
- Optionally receives normalized chat events from `gadget-plugin-chat-adapters`. _(v3)_

## v1 Functional Requirements
1. Maintain immutable point transactions and derived balances.
2. Support one-point awards via mention syntax: `@user++` and `@user ++` (space optional).
3. Support multiple recipients in one message; if a recipient appears more than once, award once.
4. Point awards are net-neutral transfers: sender balance decreases by 1, recipient balance increases by 1. Starting balance for all users is 0. Awards are rejected if the sender has insufficient balance.
5. Reject and do not write transactions for:
   - self-awards
   - DM context
   - cross-workspace context
   - deactivated sender or deactivated recipient (`deleted: true` in Slack's user object)
   - bot-originated messages
6. Include top-level messages and thread replies as valid contexts.
7. Exclude edits from award parsing and crediting.
8. `--` never removes points; trigger playful/quippy consumer response only.
9. Attempting to award the consuming bot triggers playful/quippy consumer response only.
10. Protect against double-credit with idempotency keys.

## v1 Configuration
- `feedback.reaction`: string emoji name, default `+1`; posted as a reaction to the triggering message on successful award

## Events and API Contracts
### Emitted
Emitted event schemas are deferred to v2. The event system is designed to be additive: event emission hooks will be inserted into write paths without requiring structural changes when schemas are defined.

All writes must be idempotent by `(event_id, source_id)` or equivalent dedupe key.

## Non-Goals in v1
- Slash command award route and command mode configuration (`top_level` / `subcommand`). _(v2)_
- Scheduled jobs: monthly active-user awards, weekly leaderboard publishing. _(v3)_
- Timezone boundary logic and scheduler integration. _(v3)_
- Reaction-based upvote flow. _(v2)_
- Collusion detection and enforcement. _(v2)_
- Moderator override tooling. _(v2)_
- Native event subscription model owned by this plugin. _(v2)_
- Image leaderboard rendering and dashboard deep-link generation. _(v3)_
- Transparency mode (public award announcements and open balance lookups). _(v2)_
- `gadget-plugin-spam-reports` integration (requires inter-plugin event system not yet available in Gadget). _(v2)_

## v2 Targets
- Slash command award route with configurable `top_level` and `subcommand` modes.
- Reaction-based upvotes.
- Collusion detection signals.
- Moderator override tools (reverse/freeze/exclude).
- Native event subscription model.
- Transparency mode: `transparency.enabled` config flag; bot posts a public channel message on each award and allows any member to query another member's balance.
- `gadget-plugin-spam-reports` integration: consume `spam.report.resolved` and award the first reporter on successful removal. Blocked on Gadget inter-plugin event system.

## v3 Targets
- Monthly active-user awards with configurable points and timezone boundary logic.
- Weekly per-workspace leaderboard publishing to a configurable channel.
- Active-user activity tracking via message event listener and optional `gadget-plugin-chat-adapters` source.
- Lurker exclusion heuristics beyond baseline activity checks (e.g. last-login signals via `team.accessLogs`).
- Image leaderboard publishing with dashboard link.
- Define and publish `engagement.points.awarded` event schema.
- Define and publish `engagement.leaderboard.generated` event schema.

## Extractable Issues
1. **Plugin scaffolding: `plugin.go`, config struct, `slackclient.Client` interface, `main.go` wiring**
   Milestone: `v1-core`
2. **Define ledger schema and idempotent write path**
   Milestone: `v1-core`
3. **Implement mention parser for `@user++` and `@user ++` with recipient dedupe**
   Milestone: `v1-core`
4. **Enforce eligibility rules (self/DM/cross-workspace/suspended/bot/edit exclusions)**
   Milestone: `v1-core`
5. **Add playful Penny responses for `--` and attempts to award Penny**
   Milestone: `v1-core`
6. **Track active-user activity via message event listener**
   Milestone: `v3-scheduling`
7. **Implement monthly award job with timezone boundary logic**
   Milestone: `v3-scheduling`
8. **Implement weekly per-workspace leaderboard publishing to configurable channel**
   Milestone: `v3-scheduling`
9. **Implement slash command award route with configurable `top_level` and `subcommand` modes**
   Milestone: `v2-controls`
10. **Consume `spam.report.resolved` optionally and award first reporter on successful removal**
    Milestone: `v2-controls` _(blocked on Gadget inter-plugin event system)_
11. **Investigate giver rate limiting strategy**
    Milestone: `v2-controls`
12. **Add collusion detection signals and moderator override controls**
    Milestone: `v2-controls`
13. **Add reaction-based upvotes and native event subscription model**
    Milestone: `v2-controls`

## Recommended Package Structure

### Directory / File Tree

```
gadget-plugin-engagement-ledger/
├── go.mod
├── go.sum
├── main.go                          # Wires plugin into a host bot; feature-flag optional integrations
├── plugin/
│   └── plugin.go                    # Public API surface: Register(bot) — the only symbol consumers import
├── internal/
│   ├── parser/
│   │   └── mention.go               # Parse @user++ / @user ++ from message text; returns deduplicated recipient list
│   ├── ledger/
│   │   ├── models.go                # GORM models: Transaction, Balance
│   │   ├── writer.go                # Idempotent write path; upserts Balance after each Transaction insert
│   │   └── reader.go                # Balance queries
│   ├── eligibility/
│   │   └── rules.go                 # Enforce self/DM/cross-workspace/suspended/bot/edit exclusion rules
│   ├── handlers/
│   │   ├── mention.go               # HandlerContext handler: parses message, awards points, posts feedback
│   │   └── quip.go                  # HandlerContext handler: playful responses for -- and bot-award attempts
│   └── slackclient/
│       └── client.go                # slack.Client wrapper interface for dependency injection
└── docs/
    └── specs/
        └── gadget-plugin-engagement-ledger.md
```

### Handler Signature

Handlers follow the `HandlerContext` style used in current Gadget core:

```go
func HandleMentionAward(ctx *router.HandlerContext) error {
    // ctx.Event  — the incoming Slack event
    // ctx.Client — slackclient.Client (injected; never instantiated inside the handler)
    // ctx.DB     — *gorm.DB
    // ctx.Config — plugin config struct
}
```

`router.Router`, `router.Route`, and `slack.Client` are **not** accepted as separate positional parameters. All dependencies arrive through `HandlerContext`.

### Key Design Decisions

**`internal/` sub-packages over a flat root package**
The plugin has enough distinct concerns (parsing, ledger writes, eligibility, Slack I/O) that a flat package would conflate them and make table-driven unit tests harder to scope. Sub-packages enforce clear dependency direction: `handlers` imports `parser`, `eligibility`, and `ledger`; nothing in `internal/` imports `plugin/`.

**`plugin/` as the only public API surface**
Consuming bots call `plugin.Register(bot)` and nothing else. All internal wiring (route registration, scheduler hooks, integration guards) lives inside `plugin.go`. This keeps the import surface minimal and lets internals change without breaking consumers.

**`slackclient.Client` interface**
A thin interface wrapping `*slack.Client` methods used by this plugin (e.g., `PostMessage`, `GetUsersInfo`). Handlers accept the interface, never the concrete `*slack.Client`, so tests can inject a mock without a live Slack token.

**Idempotency key formats**
| Source | Key format |
|---|---|
| Message mention award | `mention:{workspace_id}:{event_ts}:{recipient_user_id}` |
| Slash command award _(v2)_ | `cmd:{workspace_id}:{command_id}:{recipient_user_id}` |
| Spam-report award _(v2)_ | `spam_report:{workspace_id}:{report_id}:{reporter_user_id}` |
| Monthly active-user award _(v3)_ | `monthly:{workspace_id}:{year_month}:{user_id}` |

All keys are stored on the `Transaction` row; a unique index prevents duplicate inserts.

**Optional integrations via functional options _(v2/v3)_**
`main.go` will conditionally pass `plugin.With*` options to `plugin.Register(bot, cfg, opts...)` when integrations are enabled. When no options are passed the plugin compiles and runs without any optional dependency present. No options exist in v1.

### `go.mod` Dependencies

```
module github.com/gadget-bot/gadget-plugin-engagement-ledger

go 1.26

require (
    github.com/gadget-bot/gadget       v0.8.1
    github.com/slack-go/slack          v0.18.0
    gorm.io/gorm                       v1.31.1
    gorm.io/driver/mysql               v1.6.0
    gorm.io/driver/sqlite              v1.6.0   // test / local dev only
)
```

### Testing Approach

| Layer | What is tested | Tool |
|---|---|---|
| `internal/parser` | Mention regex, space variants, dedupe, no false positives | Plain `go test` table-driven |
| `internal/ledger` | Idempotent writes, balance upserts, duplicate key rejection | `go test` with SQLite in-memory |
| `internal/eligibility` | All exclusion rules, edge cases | Plain `go test` table-driven |
| Route handlers | Full request → response cycle | `gadgettest.Dispatcher` with mock `slackclient.Client` |

`gadgettest.Dispatcher` drives a handler through the full Gadget route pipeline without a live Slack connection, allowing assertion on posted messages and written transactions in a single test.
