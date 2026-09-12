---
name: discord-integration
title: Discord Integration with Hermes
description: "Use when listing Discord members, channels, messages."
version: 1.0.0
author: Hermes Agent
license: MIT
metadata:
  hermes:
    tags: [discord, gateway, members, messaging]
    related_skills: [zalo-integration]
---

# Discord Integration with Hermes

## When to Use

User asks about Discord server state (member/user list, channels, threads, recent messages) or needs Hermes Discord gateway/bot operation.

**Trigger:** Use when user asks about Discord server state — list users/members, channels, threads, recent messages — or when operating the Hermes Discord gateway/bot.

Class-level skill for Hermes + Discord. Covers gateway tools plus REST fallback when tools are insufficient.

## Tool Map (load order)

1. `discord_admin` action=`server_info` / `list_guilds` / `list_channels` — roster size, channel IDs.
2. `discord` action=`fetch_messages` / `search_members` — recent activity, prefix lookup.
3. REST fallback via `DISCORD_BOT_TOKEN` — full member enumeration (below).

Token lives in `$HERMES_HOME/.env` as `DISCORD_BOT_TOKEN`. Never print it. Source it in shell:
`source $HERMES_HOME/.env` (Windows MSYS2/bash: `source ~/AppData/Local/hermes/.env`).

## Full Member List (reliable path)

`search_members` requires a non-empty name prefix and misses most members — do NOT use it for full roster.

```bash
source $HERMES_HOME/.env
curl -s -H "Authorization: Bot $DISCORD_BOT_TOKEN" \
  "https://discord.com/api/v10/guilds/{guild_id}/members?limit=100" > members.json
```

- `limit=100` is API max per page. For >100 members paginate with `?limit=100&after={last_user_id}`.
- Response: array of `{ user: { id, username, global_name }, nick }`. Prefer `nick || global_name || username` for display.
- `server_info` gives `member_count` — verify count matches before reporting completeness.

## Channel Management (REST fallback)

`discord_admin` has no rename action. Use REST `PATCH /api/v10/channels/{id}`:

```bash
set -a; source ~/AppData/Local/hermes/.env; set +a
python3 << 'PY'
import os, json, http.client
token=os.environ['DISCORD_BOT_TOKEN']
def patch(cid, name):
    conn=http.client.HTTPSConnection("discord.com")
    body=json.dumps({"name": name}, ensure_ascii=False).encode('utf-8')
    headers={"Authorization": f"Bot {token}", "Content-Type": "application/json; charset=utf-8"}
    conn.request("PATCH", f"/api/v10/channels/{cid}", body=body, headers=headers)
    r=conn.getresponse(); data=r.read().decode('utf-8')
    j=json.loads(data)
    print(f"{r.status} {cid} -> {j.get('name','?')} msg={j.get('message','ok')[:200]}")
    conn.close()
# Repeat for each channel:
patch("CHANNEL_ID", "new-name")
PY
```

- Requires bot role with `MANAGE_CHANNELS` (bit 16) or `ADMINISTRATOR` (bit 3).
- `charset=utf-8` in Content-Type is mandatory for emoji/channel names with special chars.
- Direct curl `-d` may mangle emoji encoding; Python `http.client` with `ensure_ascii=False` is reliable.
- User convention: emoji + fullwidth dot `・` + lowercase (e.g. `💬・chat`).

## Member Kick (REST)

`discord_admin` has no kick action. Use `DELETE /api/v10/guilds/{gid}/members/{uid}` (no body) — 204 = kicked. Verify with a fresh member list.

## Partial / Alternative Paths

- Recent active users only: `discord` action=`fetch_messages` on channel_id → dedupe `author` fields. Fast but incomplete.
- Single lookup: `discord_admin` action=`member_info` with `guild_id` + `user_id`.
- Thread IDs: `$HERMES_HOME/discord_threads.json` + `channel_directory.json` map thread/channel IDs to names.

## Channel Creation (REST)

`discord_admin` has no create-channel action. Use `POST /api/v10/guilds/{guild_id}/channels`:

```bash
set -a; source ~/AppData/Local/hermes/.env; set +a
python3 << 'PY'
import os, json, http.client
token = os.environ['DISCORD_BOT_TOKEN']
def create(name, ch_type=0, parent_id=None):
    body = {"name": name, "type": ch_type}
    if parent_id: body["parent_id"] = parent_id
    conn = http.client.HTTPSConnection("discord.com")
    data = json.dumps(body, ensure_ascii=False).encode("utf-8")
    headers = {"Authorization": f"Bot {token}", "Content-Type": "application/json; charset=utf-8"}
    conn.request("POST", f"/api/v10/guilds/{guild_id}/channels", data, headers)
    r = conn.getresponse(); resp = r.read().decode()
    j = json.loads(resp) if resp else {}
    print(f"{r.status} created: {j.get('name','')} id={j.get('id','')} msg={j.get('message','')[:200]}")
    conn.close(); return j
PY
```

- `type`: 0=text, 2=voice, 4=category, 5=announcement, 15=forum.
- Optional: `topic`, `parent_id` (for category placement), `nsfw`, `rate_limit_per_user`.
- Use Python `http.client` with `ensure_ascii=False` for emoji-safe names.

## Cron Job Delivery to Discord

When creating cron jobs, use `deliver: "discord:<channel_id>"` to send output to a specific Discord channel.

```python
cronjob_manage(action='create', 
    schedule='15m',
    deliver='discord:<channel_id>',
    prompt='...')
```

- `discord:<channel_id>` — sends to channel main chat.
- `discord:<channel_id>:<thread_id>` — sends to a specific thread.
- `origin` — sends back to the chat where the cron was created.
- `all` — every connected home channel.
- Recurring schedule forms: `15m`, `every 2h`, `*/15 16-23 * * *` (cron syntax).
- One-shot: `in 30m`, or ISO timestamp.
- `attach_to_session: true` makes output continuable (thread reply context).
- Always verify delivery target exists before creating the job.
- **Brevity constraint**: Keep cron prompts delivering to Discord strictly bounded (e.g. max 3-5 items, exactly 1 line per bullet `• [Short title]: [URL]`, no preamble or commentary) to prevent wall-of-text spam in chat channels.

## Pitfalls

| Failure | Cause | Fix |
|---|---|---|
| `search_members` returns 0 with query `a` / empty query errors | Tool is prefix-search, not list-all | Use REST `GET /guilds/{id}/members` fallback |
| `/tmp/*.json` not found (Windows MSYS2) | `/tmp` path mismatch in tool sandbox | Write to CWD (`./members.json`), not `/tmp/` |
| f-string quoting breaks under `source .env` + bash eval | Nested quotes through execute_code/terminal layers | Keep inline python minimal; prefer simple prints, avoid nested f-string quotes |
| `os.environ` empty after `source .env` | `source` sets shell vars but doesn't export by default | Use `set -a; source .env; set +a` before Python |
| `curl PATCH` with emoji → `invalid JSON` | Shell quoting mangles UTF-8/emoji in `-d` | Use `http.client` + `json.dumps(..., ensure_ascii=False).encode('utf-8')` |
| `DELETE /channels/{category}` → 400 category not empty | Discord requires empty category | Move/delete child channels first, then delete category |
| `DELETE /channels/{id}` → 404 Unknown Channel | Already deleted / stale ID from cached `list_channels` | Ignore 404, re-verify with fresh `list_channels` |
| Cron job created but no output visible in channel | Wrong channel_id or channel locked for bot | Verify channel permissions allow bot SEND_MESSAGES; check `deliver` format is `discord:<id>` (not `discord:<name>`) |
| Cron schedule `15m` fires too often / outside hours | User wanted specific hours | Use cron syntax: `*/15 16-23 * * *` (every 15min, 4pm-11pm only) |
| `error code: 1010` / Cloudflare block on REST `curl` | WARP + 9Router MITM + proxy env vars intercept Discord TLS → flagged IP | Bypass proxy: `env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY curl -s -X POST ... -H "User-Agent: DiscordBot (OX, 1.0)"`; never rely on default `https_proxy` |
| Single POST with 20 `@mentions` → 1010 / spam flag | Discord anti-spam rate-limits bulk mentions in one message | Batch 5 mentions per message, 1.5s delay between batches, `shlex.quote(json.dumps(payload))` for safe quoting; send all batches to same welcome/channel ID |
| Cron output tables don't render on Discord | Discord markdown has no table support | Cron prompt must produce bullet lists with inline links (`- **Category** — summary [source](url)`) not markdown tables |

## Community Reset Recipe (3 categories)

Standard layout: `📋・information` (read-only), `💬・general` (open), `🔊・voice`.

```bash
set -a; source ~/AppData/Local/hermes/.env; set +a
python3 << 'PY'
import os, json, http.client
token=os.environ['DISCORD_BOT_TOKEN']
def req(m, p, b=None):
    import http.client, json
    c=http.client.HTTPSConnection("discord.com")
    d=json.dumps(b, ensure_ascii=False).encode('utf-8') if b else None
    h={"Authorization": f"Bot {token}", "Content-Type": "application/json; charset=utf-8"}
    c.request(m, p, body=d, headers=h); r=c.getresponse(); t=r.read().decode('utf-8')
    import json as _j; j=_j.loads(t) if t else {}
    print(m, p, r.status, j.get('name',''), j.get('message','')[:200]); c.close(); return r.status, j
# category type=4
req("POST","/api/v10/guilds/{gid}/channels",{"name":"📋・information","type":4})
# text type=0 + parent_id, voice type=2 + parent_id
req("POST","/api/v10/guilds/{gid}/channels",{"name":"📜・rules","type":0,"parent_id":cat_id})
# move: PATCH /channels/{id} {"parent_id": new_cat}
# lock read-only: PUT /channels/{ch}/permissions/{everyone} {"id":role,"type":0,"allow":"1024","deny":"2048"} # VIEW=1024 SEND=2048, 204=ok
# delete: DELETE /channels/{id} (ignore 404)
PY
```

- Naming: `{emoji}・{kebab-lowercase}` (U+30FB `・`), not `|` or `-`.
- Keep home channel (`DISCORD_HOME_CHANNEL` 1152202940724019281 `💬・chat`): move via `PATCH parent_id` instead of recreating.
- Verify with `discord_admin` action=`list_channels` after batch.
