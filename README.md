# OxHUD - Universal Agent HUD & Statusline

Zero-dependency, high-performance Go statusline for:
- **Claude Code** (`.claude-plugin`, statusline API)
- **Hermes Agent** (native skill & hooks)
- **OpenAI Codex** (CLI statusline adapter)

## Quick Start

```bash
# Build
go build -o bin/oxhud.exe ./cmd/oxhud

# Test Claude mode
echo '{"model":{"display_name":"Opus 3.7"},"context_window":{"current_usage":{"input_tokens":45000,"output_tokens":5000},"context_window_size":200000}}' | ./bin/oxhud.exe --agent=claude

# Test Hermes mode
echo '{"model":"anti","used_tokens":85000,"total_tokens":256000,"active_tool":"terminal","task":"executing build"}' | ./bin/oxhud.exe --agent=hermes

# Test Codex mode
echo '{"model":"gpt-4o","used_tokens":110000,"total_tokens":128000,"active_tool":"edit_file","status":"patching"}' | ./bin/oxhud.exe --agent=codex
```

## Plugins & Skills Structure
- `plugins/claude/`: Claude Code plugin manifest & setup command
- `skills/hermes/`: Hermes Agent skill spec (`SKILL.md`)
- `skills/codex/`: Codex pipeline integration guide
