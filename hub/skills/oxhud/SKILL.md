---
name: oxhud
description: Universal real-time CLI HUD & Statusline for Hermes, Claude Code, and Codex.
author: TheOwlOps
version: 0.1.0
---

# OxHUD - Universal Agent HUD

Use `oxhud` to inspect real-time agent token usage, context window consumption, active tools, and background subagent tasks.

## Commands

- **Render statusline directly**:
  ```bash
  oxhud --agent=hermes
  oxhud --agent=claude
  oxhud --agent=codex
  ```

- **Pipe JSON state**:
  ```bash
  echo '{"model":"anti","used_tokens":50000,"total_tokens":256000,"active_tool":"terminal","task":"running tests"}' | oxhud --agent=hermes
  ```

## Integration for Hermes

Add or invoke `oxhud` in post-tool hooks or status reporting workflows to show compact single-line terminal feedback.
