# gadget-plugin-engagement-ledger Specification

## Purpose
Provide a standalone engagement points system with an immutable ledger, user-to-user transfers, monthly activity awards, and per-workspace leaderboards.

## Standalone and Optional Integrations
- Operates standalone with direct command/API input.
- Optionally consumes events from `gadget-plugin-spam-reports`.
- Optionally registers recurring jobs with `gadget-plugin-scheduler`.
- Optionally receives normalized chat events from `gadget-plugin-chat-adapters`.

## v1 Functional Requirements
1. Maintain immutable point transactions and derived balances.
2. Support one-point awards via:
   - mention syntax: `@user++` and `@user ++` (space optional)
   - slash command route (prefix-configurable)
3. Support multiple recipients in one message; if a recipient appears more than once, award once.
4. Net-neutral transfers only (sender decreases, recipient increases).
5. Reject and do not write transactions for:
   - self-awards
   - DM context
   - cross-workspace context
   - suspended sender or suspended recipient
   - bot-originated messages
6. Include top-level messages and thread replies as valid contexts.
7. Exclude edits from award parsing and crediting.
8. `--` never removes points; trigger playful/quippy consumer response only.
9. Attempting to award the consuming bot triggers playful/quippy consumer response only.
10. Protect against double-credit with idempotency keys.
11. Run monthly active-user awards:
    - default `10` points per active user, configurable by plugin consumer
    - active user = posted top-level message or thread reply in channels in that month where the bot is a member
    - inactive users receive nothing
    - timezone boundary mode configurable per workspace: `workspace_local` or `utc`
12. Publish weekly leaderboard per workspace to a configurable channel.

## v1 Configuration
- `commands.mode`: `top_level | subcommand`
- `commands.prefix`: string (known prefix namespace)
- `monthly_award.enabled`: boolean
- `monthly_award.points`: integer, default `10`
- `monthly_award.timezone_mode`: `workspace_local | utc`
- `leaderboard.enabled`: boolean
- `leaderboard.channel_id`: string
- `transparency.enabled`: boolean
- `integrations.spam_reports.enabled`: boolean
- `integrations.scheduler.enabled`: boolean
- `integrations.chat_adapters.enabled`: boolean

## Events and API Contracts
### Consumed (optional)
- `spam.report.resolved`
  - Required fields: `event_id`, `workspace_id`, `report_id`, `removed`, `first_reporter_user_id`
  - Behavior: if `removed=true`, award first reporter exactly once.

### Emitted
- `engagement.points.awarded`
- `engagement.leaderboard.generated`

All writes must be idempotent by `(event_id, source_id)` or equivalent dedupe key.

## Non-Goals in v1
- Reaction-based upvote flow.
- Collusion detection and enforcement.
- Moderator override tooling.
- Native event subscription model owned by this plugin.
- Image leaderboard rendering and dashboard deep-link generation (text publish first; add when UI support exists).

## v2 Targets
- Reaction-based upvotes.
- Collusion detection signals.
- Moderator override tools (reverse/freeze/exclude).
- Native event subscription model.
- Lurker exclusion heuristics beyond baseline activity checks.
- Image leaderboard publishing with dashboard link.

## Extractable Issues
1. **Define ledger schema and idempotent write path**  
   Milestone: `v1-core-engagement`  
   Labels: `type:feature`, `area:engagement`, `priority:p0`, `standalone`
2. **Implement mention parser for `@user++` and `@user ++` with recipient dedupe**  
   Milestone: `v1-core-engagement`  
   Labels: `type:feature`, `area:engagement`, `priority:p0`
3. **Implement command routing with configurable `top_level` and `subcommand` modes**  
   Milestone: `v1-core-engagement`  
   Labels: `type:feature`, `area:engagement`, `type:api`, `priority:p0`
4. **Enforce eligibility rules (self/DM/cross-workspace/suspended/bot/edit exclusions)**  
   Milestone: `v1-core-engagement`  
   Labels: `type:feature`, `area:engagement`, `priority:p0`
5. **Add playful Penny responses for `--` and attempts to award Penny**  
   Milestone: `v1-core-engagement`  
   Labels: `type:feature`, `area:engagement`, `priority:p1`
6. **Implement monthly active-user awards with configurable points and timezone mode**  
   Milestone: `v1-core-engagement`  
   Labels: `type:feature`, `area:engagement`, `priority:p0`
7. **Implement weekly per-workspace leaderboard publishing to configurable channel**  
   Milestone: `v1-optional-integrations`  
   Labels: `type:feature`, `area:engagement`, `priority:p1`, `integration:optional`
8. **Consume `spam.report.resolved` optionally and award first reporter on successful removal**  
   Milestone: `v1-optional-integrations`  
   Labels: `type:feature`, `area:engagement`, `area:spam`, `integration:optional`, `priority:p0`
9. **Investigate giver rate limiting strategy**  
   Milestone: `v2-advanced-controls`  
   Labels: `type:feature`, `area:engagement`, `priority:p1`
10. **Add collusion detection signals and moderator override controls**  
    Milestone: `v2-advanced-controls`  
    Labels: `type:feature`, `area:engagement`, `priority:p1`
11. **Add reaction-based upvotes and native event subscription model**  
    Milestone: `v2-advanced-controls`  
    Labels: `type:feature`, `area:engagement`, `priority:p1`
