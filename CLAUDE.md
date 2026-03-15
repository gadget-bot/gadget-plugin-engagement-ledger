# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`gadget-plugin-engagement-ledger` provides a standalone engagement points system with an immutable ledger, user-to-user transfers, monthly activity awards, and per-workspace leaderboards. It optionally consumes events from `gadget-plugin-spam-reports`, registers recurring jobs with `gadget-plugin-scheduler`, and receives normalized chat events from `gadget-plugin-chat-adapters` — but has no hard dependency on any of them.

This plugin is part of the [Gadget](https://github.com/gadget-bot/gadget/) ecosystem — a Go Slack bot framework built on the Slack Events API.

See `docs/specs/gadget-plugin-engagement-ledger.md` for full v1 requirements.

## Build & Development Commands

```bash
make build        # Compile binary to dist/
make test         # Run tests with coverage
make lint         # Run golangci-lint (also runs fmt check)
make fmt          # Check formatting with golangci-lint (diff only)
make all          # clean → verify → lint → test → build
make tools        # Install golangci-lint
make start-db     # Start local MariaDB container
make stop-db      # Stop local MariaDB container
make container    # Build Docker image as gadget-plugin-engagement-ledger:local
```

Run a single test:
```bash
go test -v -run TestFunctionName ./pkg/...
```

**Always prefer `make` targets over calling `go` commands directly.** Only fall back to raw `go` commands when no suitable make target exists (e.g., running a single test).

Run `make build` and `make test` after any code changes before committing.

## Architecture

<!-- TODO: document key packages, entry points, and data flow as the plugin grows. -->

## Testing

Tests use Go's table-driven testing pattern. Always use interfaces and dependency injection so handlers accept mock clients rather than creating their own. Never instantiate real API clients inside handlers.

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`

**Breaking changes:** append `!` after type/scope, or add `BREAKING CHANGE:` in the footer.

## Branching

- `main` — production releases only; NEVER commit directly to main
- Always create a feature branch and open a PR for all changes

## GitHub Repository

The origin repository is `gadget-bot/gadget-plugin-engagement-ledger`. Always use this owner/repo when querying GitHub for issues, pull requests, or other repository details.

## GitHub Issues

When opening issues:
- Always apply an issue type and the best-fitting label
- Scan existing issues to identify relationships (sub-issues, duplicates, related issues)
- Ask before changing existing relationships
