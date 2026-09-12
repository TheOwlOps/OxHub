# Hermes 1:1 Architecture Blueprint for OX Alpha

Mapping Hermes Agent (`nousresearch/hermes-agent`) domain layout to the Tri-Stack architecture (Go Core + Python Compute + TypeScript Presentation).

## Directory Layout Comparison

| Hermes Agent Path | OX Alpha Tri-Stack Path | Tech Stack | Responsibility |
| :--- | :--- | :--- | :--- |
| `agent/` | `agent/` | Go | ReAct Turn Loop, context compaction, prompt assembly, mid-turn steering |
| `tools/` | `tools/` | Go + Python | Core filesystem/shell tools in Go; CDP browser & sandboxed exec in Python worker |
| `gateway/` | `gateway/` | Go | Multi-platform messaging hub (`platforms/`: Telegram, Discord, Zalo, Webhook) |
| `cron/` | `cron/` | Go | Native goroutine scheduler (`robfig/cron/v3`) & delivery queue |
| `providers/` | `providers/` | Go | LLM provider runtime resolution (OpenAI-compatible, Claude, fallback chains) |
| `skills/` | `skills/` | Go + MD | Procedural skill scanner, `skill_view`, dynamic `skill_manage` ops, `user-skills/` |
| `state/` | `state/` | Go | SQLite WAL database (`ox.db`), session persistence, FTS message search |
| `web/` | `web/` | TS / React | Web SPA dashboard, Cowork Canvas, ReAct Console (Pure Presentation) |
| `apps/desktop/` | `apps/desktop-wails/` | Go + TS | Wails v2 desktop shell embedding `web/` bundle |

## Technology Boundaries

1. **Go (System Core)**:
   - Zero-dependency compilation (`CGO_ENABLED=0`) with `modernc.org/sqlite`.
   - Native Goroutines for multi-agent delegation worker pool and cron scheduler.
   - Long-lived Gateway connections: WebSocket for Discord, long-polling for Telegram, SSE client for Zalo.

2. **Python (Heavy Compute & AI Ecosystem)**:
   - Isolated environment managed via `uv`.
   - Headless Chrome CDP control via Playwright for `browser_exec`.
   - Code execution sandbox (`execute_code`) with process timeouts and scratch cleanup.
   - Heavy document parsing (PDF, DOCX, OCR).
   - Go communicates with Python worker via Stdio JSON-RPC (zero network port attack surface).

3. **TypeScript / React (Presentation)**:
   - Stateless UI interacting strictly via REST/WebSocket API (`http://localhost:1060`).
   - Real-time token stream rendering, biological mascot animation states, and Cowork workspace.
