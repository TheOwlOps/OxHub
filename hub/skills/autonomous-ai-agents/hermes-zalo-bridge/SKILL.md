---
name: hermes-zalo-bridge
title: Hermes ⟷ Zalo concierge bridge (zca-js)
description: "Use when operating Zalo gateway. No code in replies."
version: 0.1.0
author: Hermes Agent
license: MIT
platforms: [windows]
tags: [zalo, zca-js, gateway, bridge, concierge, no-code-replies]
related_skills: [hermes-agent]
---

# Hermes ⟷ Zalo Concierge Bridge

Runs Hermes Agent through a Zalo personal account (`37082029033239591`) so messages
inside a Zalo group land as gateway turns. Unofficial integration via
`zca-js`; Zalo may invalidate sessions without warning.

### Quick start

```bash
# 1. Bridge (Node.js) must be running on localhost:8787
pm2 start server.js --name monitor           # from the plugin dir
pm2 save                                     # persist to PM2 dump
pm2 resurrect                                # restore after reboot

# 2. Health check
curl -s http://127.0.0.1:8787/health
# expect: {"ok":true,"loggedIn":true,"sessionDead":false,...}

# 3. Gateway connects automatically once bridge is healthy.
hermes gateway restart                       # if SSE stays at 0 clients
```

## Tone Enforcement (user-mandated)

Zalo replies MUST follow these rules, derived from SOUL.md and user corrections:
- **Prefix**: start with **dạ / vâng / ạ** depending on context
- **Full sentences only** — no telegraphic/terse fragments
- **Length**: 2–4 lines, warm and respectful
- **Ending**: always end with **ạ** or **dạ**
- **Auto-@**: every group reply MUST include `@member` mention (pos=0, len=7) for the asker
- **No code blocks** — prose/list/link only; inline backticks OK, fenced blocks forbidden
- **Think first**: verify facts before replying; no hallucination

Discord/Telegram stay ultra-terse; Zalo is polite full-sentence mode — never mix them.

## Mention + welcome capture

- Every `@bot` in the allowlisted group stores the mentioned UID in
  `~/.hermes-zalo/mentioned_uids.json`.
- **Native group_event welcome (preferred):** `adapter.py` handles `group_event`
  SSE events from zca-js. On `type == "join"`, it reads `data.data.groupId` for
  thread_id, `data.data.groupName` for dynamic group name, and
  `data.data.updateMembers[].id` for UIDs. Sends a multi-line welcome with
  emoji and `@member` mentions (pos=0, len=7 per UID, spaced 8 chars apart).
- **Critical:** `threadId` is NOT at top level of group_event data — use
  `data.get("data", {}).get("groupId")` as fallback. Original code that read
  `data.get("threadId")` always returned empty string.
- **Critical:** `groupName` must be read from `data.data.groupName`, never
  hardcoded. When group is renamed, the event carries the new name.
- **Auto-@ dynamic reply (100 users):** `_on_inbound_message` stores
  `senderId` per `threadId` in `self._last_sender[threadId]`. Override
  `async def send(...)` to inject `mentions=[{"pos":0,"uid":senderId,"len":7}]`
  and prepend `@member ` when replying in group. UID is per-message, not
  fixed — handles 100 concurrent askers. Fallback to plain reply if missing.
- @reply format: `@member <response>` with blue mention via `mentions` object.

## Login / QR flow

`node login.mjs --force` writes `qr.png` to `~/.hermes-zalo/`.
Scan with Zalo mobile app once. `credentials.json` persists across reboots.

## Recovery checklist

| Symptom | Fix |
|---|---|
| `sessionDead: true` / `loggedIn: false` | `node login.mjs --force`, re-scan QR |
| Bridge 403 / EADDRINUSE | `pm2 restart monitor`; kill stale `node` PIDs first |
| SSE clients stuck at 0 | `hermes gateway restart` |
| Gateway not in config list | Re-enable plugin: `hermes plugins enable zalo-platform --allow-tool-override` |

## Display suppression (anti-spam for Zalo)

Zalo must receive only the final answer — no interim/stream spam. Set in `C:\Users\OS\AppData\Local\hermes\config.yaml` and restart gateway **from outside** the gateway (separate shell / `pm2 restart`, not `hermes gateway restart` from inside):

```yaml
display:
  interim_assistant_messages: false
  tool_progress: null        # not "none" — null disables
  long_running_notifications: false
  busy_ack_detail: false
  streaming: false
  memory_notifications: 'off'  # hides "Self-improvement review: Memory updated"
```

Verify with `cat config.yaml` after restart; `streaming: false` prevents `Working — X min — iteration Y/150` and `waiting on research` frames from ever reaching Zalo.

## Pitfalls

- **Double @mention when streaming is on:** `adapter.py` injects `@member` on `idx==0` chunk. With streaming, both streamed chunks and final `send()` fire → two blue tags in one answer. Fix: `display.streaming: false` → single final send only.
- **Bad JSON kills the bridge (fixed 2026-09-09):** zca-js can emit a message
  whose body has a bad escaped char, and `express.json()` throws
  `SyntaxError: Bad escaped character in JSON at position N` → unhandled → the
  whole PM2 `monitor` process crashes → port 8787 refused → gateway can't
  connect. **Fix:** add an express JSON error-handler middleware right after
  `app.use(express.json({ limit: "2mb" }))`:
  ```js
  app.use((err, req, res, next) => {
    if (err instanceof SyntaxError && err.status === 400 && "body" in err) {
      console.error("[bridge] bad JSON body skipped:", err.message?.slice(0, 200));
      return res.status(400).json({ error: "bad JSON", detail: err.message });
    }
    return next(err);
  });
  ```
  This drops the bad body with a 400 + log and keeps the bridge alive. Also
  guard `client.on("message")` / poller bodies in try/catch (the welcome
  poller already does). Symptom that precedes it: `[bridge] bad JSON` /
  `[welcome-poller] check failed: Retry limit` in `monitor-error.log`, then
  `ConnectionRefusedError` port 8787 in gateway.log.
- **Duplicate connections**: never run two Zalo sessions; Zalo kicks the
  older one (`3003 KICKOUT_BY_WORKER`). One PM2 process only.
- **Windows CRLF**: use `\n` in JS strings sent over bridge; `\r\n` breaks SSE.
- **Cookie expiry**: auto-relagins runs 5× backoff on `code=1006`. If it keeps
  failing, delete `credentials.json` and re-scan QR.

## References

- `references/zalo-troubleshooting.md` — full error transcript history
- `references/zalo-config.yaml` — sample gateway config block
