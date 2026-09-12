# Codex Integration for OxHUD

Use `oxhud` inside Codex execution pipelines or CLI prompt status wrapper.

## Usage

In any bash / node sub-process spawned by Codex:

```bash
echo "{\"model\":\"$CODEX_MODEL\",\"used_tokens\":$USED_TOKENS,\"total_tokens\":128000,\"active_tool\":\"$ACTIVE_TOOL\",\"status\":\"$STATUS\"}" | oxhud --agent=codex
```

## Wrapper Helper

A tiny wrapper script `codex-hud.sh` or `codex-hud.bat`:

```cmd
@echo off
oxhud --agent=codex %*
```
