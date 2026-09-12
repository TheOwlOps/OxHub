---
name: gateway-troubleshooting
description: "Diagnose Hermes Gateway crashes and model errors."
version: 1.1.0
author: Hermes Agent
license: MIT
platforms: [linux, macos, windows]
metadata:
  hermes:
    tags: [gateway, troubleshooting, telegram, discord, crash, oom, auto-start]
    related_skills: [hermes-agent]
---

# Gateway Troubleshooting

Diagnose and fix Hermes Gateway issues. Gateway is the background process that connects Hermes to messaging platforms (Telegram, Discord, Slack, etc.) and runs cron jobs.

## When to Use

- When the gateway process is down or not responding to messages
- When Telegram/Discord/Slack bot stops replying but shows "connected"
- When gateway crashes repeatedly with SIGKILL/OOM patterns
- When model validation errors block the gateway from starting
- When setting up auto-start for gateway on Windows or Linux

## Quick Status Check

```bash
# Is gateway running?
hermes gateway status

# Check recent logs
tail -50 ~/AppData/Local/hermes/logs/gateway.log        # Windows
tail -50 ~/.hermes/logs/gateway.log                      # Linux/macOS

# Check for errors specifically
grep -i "ERROR|error|crash|SIGKILL|OOM" ~/AppData/Local/hermes/logs/gateway.log | tail -20
```

## Common Failure Patterns

### 1. Gateway Not Running
- **Symptom:** `hermes gateway status` shows "No gateway process detected"
- **Fix:** `hermes gateway start`
- **If it starts then immediately dies:** Check logs for model/config errors (see #3 below)

### 2. Gateway Crashes Repeatedly (SIGKILL / OOM)
- **Symptom:** Logs show `exited UNCLEANLY (no exit path ran — SIGKILL / OOM / VM death)`
- **Cause:** Process killed by OS (out of memory) or forced termination
- **Fix:** Restart with `hermes gateway start`. If recurring, check system memory usage.
- **Note:** On Windows, check if another process (antivirus, system cleanup) is killing it.

### 3. Model Context Length Validation Error
- **Symptom:** `ValueError: Model X has a context window of N tokens, which is below the minimum 64,000`
- **Cause:** A custom model in `config.yaml` has `context_length` set too low (e.g., 1)
- **Fix:**
  ```bash
  hermes config set custom_providers.0.models.<model_name>.context_length 256000
  ```
- **Why it happens:** The local model server does not report context_length, so Hermes falls back to the config value. If that is wrong (e.g., 1), the agent refuses to start.

### 4. Platform Connected but No Response
- **Symptom:** Logs show "Connected to Telegram/Discord" but messages get no reply
- **Check:**
  1. Is there an `Agent error` in the logs after the inbound message?
  2. Is the model endpoint reachable? (`curl http://localhost:<port>/v1/models`)
  3. Are there empty-content transcript errors? (`Pre-call sanitizer: healed empty messages`)
  4. Is the upstream model quota exhausted (HTTP 429 `RESOURCE_EXHAUSTED` / `Individual quota reached`)? Check proxy logs/response error strings.
- **Fix:** Restart gateway after fixing the underlying issue or wait/swap upstream credential pool

### 6. Upstream Proxy / Local Proxy Model Troubleshooting (e.g., 9Router / Port 20128)
- **429 RESOURCE_EXHAUSTED / Quota Exhausted**: Account quota reached on OAuth provider. Injecting new refresh tokens into SQLite DB requires restarting the proxy process to clear RAM locks/backoff timers (`modelLock_*`).
- **403 VALIDATION_REQUIRED (Google Anti-abuse Checkpoint)**: Google Cloud Code AI backend locks account requiring SMS verification to shortcode `96831`. Non-US SIM carriers (Vietnam Viettel/Vina/Mobi, etc.) block outbound SMS to US shortcodes, creating an unresolvable verification loop. Such accounts cannot be used for Antigravity until checkpoint lifts or verified with US carrier.
- **403 Organization Restriction**: Workspace/edu emails (`@*.edu.vn`, etc.) restricted by tenant admin: `"You're restricted from using Gemini Code Assist for individuals in your organization."` Must use personal Gmail account.
- **404 Requested entity was not found**: Upstream model not supported or project mapping mismatch on health-check; usually benign if fallback models (e.g. `gemini-3.8-flash`) respond 200.
- **405 Method Not Allowed**: Upstream provider is temporarily resetting or rate-limiting (e.g., 30s reset window); retry after the window expires.
- **400 Bad Request (Codex account)**: Specified model is not supported with ChatGPT Codex accounts.

### 5. Platform Never Connects
- **Telegram:** Check `TELEGRAM_BOT_TOKEN` in `.env` (not commented out with `#`)
- **Discord:** See "Connecting Discord from scratch" below.
- **All platforms:** Verify network connectivity to the platform API

## Connecting Discord from scratch (Windows)

Token setup has 3 traps — worth doing right the first time.

1. **Token MUST be env var, not a config.yaml key.** Setting
   `hermes config set gateway.discord.bot_token X` writes a top-level key that
   Hermes *saves but ignores* ("not a recognized config key"). The adapter
   only reads the env var. Put it in `~/.hermes/.env` (`AppData/Local/hermes/.env`):
   ```
   DISCORD_BOT_TOKEN=...
   DISCORD_ALLOW_ALL_USERS=true
   ```
   `DISCORD_ALLOW_ALL_USERS=true` is REQUIRED or every Discord message is
   denied: log line `[Discord] Discord messages are being denied because no
   allowlist is configured`. Without it the bot shows online but silently
   ignores everyone.
2. **Enable MESSAGE CONTENT INTENT** in Developer Portal → Bot →
   Privileged Gateway Intents. Without it the bot can't read message content.
   After toggling, restart gateway so the token+intents reload.
3. **`Improper token has been passed` / `401 Unauthorized`** in gateway.log =
   the token in `.env` is no longer valid. Almost always because the token was
   leaked/shared publicly (Discord auto-revokes) or reset on the portal. Fix:
   Portal → Bot → Reset Token → copy new → **replace the single line** in
   `.env` (editing via `sed -i '/DISCORD_BOT_TOKEN/d'` then re-append avoids
   duplicate `DISCORD_BOT_TOKEN=` lines) → restart gateway.

Invite URL from Portal → OAuth2 → URL Generator: tick `bot` +
`applications.commands` + perms (Send Messages, Read Message History, Embed
Links, Use Slash Commands). Do NOT tick `identify`/`email`/`guilds` scopes or
`Administrator` — extra scopes trigger "Integration requires code grant" /
"high-risk" errors. If the error persists on a specific app, the app has a
stuck OAuth code-grant flag → create a fresh Application (fastest).

Confirmation the adapter is live: `[Discord] Connected as NAME#TAG` then
`✓ discord connected`. If you only see `Connecting to discord...` looping with
`LoginFailure`, it's the token.

## Windows Auto-Start

Two approaches for keeping gateway alive across reboots:

### Option A: Task Scheduler (runs after user login)
```bash
hermes gateway install    # registers scheduled task "Hermes_Gateway"
```
- Runs when user logs in
- Suitable for most users
- Verify with: `hermes gateway status` shows "Scheduled Task registered: Hermes_Gateway"

### Option B: Windows Service (runs at boot, before login)
- More reliable for 24/7 operation
- Requires NSSM or similar service wrapper
- Gateway runs even without user session

### Restarting after config changes
After editing `.env` or `config.yaml`:
```bash
hermes gateway restart
```

## Log Analysis

Key log files (Windows path: `C:\Users\OS\AppData\Local\hermes\logs\`):

- `gateway.log` — main gateway lifecycle and platform events
- `gateway-stdio.log` — full stderr/stdout including tracebacks
- `gateway-error.log` — error-level messages
- `gateway-exit-diag.log` — diagnostics from unclean exits

### Key log patterns to watch for

- **SIGKILL / OOM**: `exited UNCLEANLY (no exit path ran — SIGKILL / OOM / VM death)` — gateway killed by OS
- **Model error**: `ValueError: Model X has a context window of N tokens` — fix with `hermes config set`
- **Agent error**: `ERROR gateway.run: Agent error in session agent:main:telegram:dm:...` — agent failed to start
- **Connection success**: `✓ telegram connected` followed by `Gateway running with 1 platform(s)`
- **Transcript warning**: `Pre-call sanitizer: healed 1 empty non-final message(s)` — self-recovering, usually harmless

### Diagnostic flow when bot stops responding

1. `hermes gateway status` — is it running?
2. `grep -i "inbound message" gateway.log | tail -10` — recent messages received
3. `grep -A5 "Agent error" gateway-stdio.log | tail -30` — check for agent failures
4. `curl -s http://localhost:20128/v1/models | head -5` — verify model endpoint
5. `grep TELEGRAM_BOT_TOKEN ~/.hermes/.env` — confirm token is uncommented
6. `hermes gateway restart` — restart after fixes
7. Send test message, re-check logs

## Windows Auto-Start Detailed

### Task Scheduler
- Status check: `hermes gateway status` shows "Scheduled Task registered"
- If `Last Run Result: 1` — previous instance failed, check logs first

### Service wrapper (NSSM)
- Install: `nssm install HermesGateway <path-to-hermes> gateway run`
- Set auto-start: `nssm set HermesGateway Start SERVICE_AUTO_START`
- Start: `net start HermesGateway`

## References

See `references/gateway-windows-autostart.md` and `references/gateway-log-analysis.md` for detailed information.
