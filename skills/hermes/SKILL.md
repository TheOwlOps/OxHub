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
  D:/OxHub/bin/oxhud.exe --agent=hermes
  D:/OxHub/bin/oxhud.exe --agent=claude
  D:/OxHub/bin/oxhud.exe --agent=codex
  ```

- **Pipe JSON state**:
  ```bash
  echo '{"model":"anti","used_tokens":50000,"total_tokens":256000,"active_tool":"terminal","task":"running tests"}' | D:/OxHub/bin/oxhud.exe --agent=hermes
  ```
