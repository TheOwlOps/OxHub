---
name: zalo-integration
title: Zalo Integration with Hermes
description: "Use when connecting Hermes to Zalo messaging via bridge."
category: integrations
---

# Zalo Integration with Hermes — Skill Library

**Trigger:** Use when integrating Hermes with Zalo messaging platform; replicate OpenClaw approach; fix ZaloClient.js execution errors (running JS in cmd instead of Node.js); configure group mention-trigger mode, blue @All mentions, or new-member welcome.

**Class-level skill** for Hermes + Zalo integration. Covers bridge setup, mention tracking, welcome-on-join polling, and allowed-group restriction.

---

## 🛠️ Setup

### 1. Install Zalo plugin
```bash
hermes plugins install cuongdev/hermes-zalo-plugin --enable
```

### 2. Install PM2 + auto-start
```bash
npm install -g pm2
pm2 start server.js --name zalo-bridge
pm2 save
# Create Windows Startup batch in C:\Users\OS\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\start_zalo.bat:
@echo off
cd /d "C:\Users\OS\AppData\Local\hermes\plugins\hermes-zalo-plugin"
call "C:\Users\OS\AppData\Local\hermes\plugins\hermes-zalo-plugin\node_modules\.bin\pm2.cmd" resurrect
```

### 3. Configure Hermes gateway
```bash
hermes config set zalo.allowed_threads 8813682027038154228
hermes config set zalo.allowed_groups 8813682027038154228
hermes plugins enable zalo-platform
```

### 4. Restart gateway
```bash
hermes gateway restart
```

---

## 📋 Mention + Welcome Tracking

### a) `@mention` UID Capture

When anyone `@` another user in a Hermes Zalo group, the bridge's `server.js` patch captures the mentioned UID and appends it to `mentioned_uids.json` under `~/.hermes-zalo/`:

```json
{"uids": ["37082029033239591", "1234567890123456789"], "last_updated": "2026-08-29T..."}
```

**Purpose**: Build an allowlist of group members for future targeted replies.

### b) New-Member Welcome Polling

Since Zalo's zca-js does not emit a native `welcome` event on join, the bridge runs a **15-second polling loop** that:

1. Calls `getGroupInfo(groupId)` to fetch current `memberIds`
2. Compares against the cached list from the previous tick
3. Sends a welcome message to any **new** member UIDs:
   ```
   Chào mừng <@UID> đã đến với nhóm Hermes! 🎉
   @UID để hỏi bot nha.
   ```

**Cache**: `group_members_cache.json` stores the last-seen member list per group ID.

---

## 📁 Support Files (under `~/.hermes-zalo/`)

| File | Purpose |
|---|---|
| `mentioned_uids.json` | Auto-populated UIDs of users @mentioned in Hermes groups |
| `group_members_cache.json` | Last-seen member list for welcome-poller (per-group) |

---

## ⚠️ Pitfalls (Learned This Session)

| Failure | Cause | Fix |
|---|---|---|
| **Restart wrong process** | PM2 task was named `monitor` (from `npm start` script), not `zalo-bridge`. `pm2 restart zalo-bridge` silently failed (`Process not found`), leaving old broken process alive. | Run `pm2 list` first → use the **actual** process name in restart. Rename via `pm2 start ... --name zalo-bridge` to match expectation. |
| **`pm2 startup` fails on Windows** | PM2's autostart hooks target systemd/init.d (POSIX). On Windows, `pm2 startup` throws `Init system not found`. | Create a **`.bat` file in the Windows Startup folder** that runs `pm2 resurrect`: see Setup step 2. |
| **`getGroupInfo()` returns `memberIds: []`** | Zalo's `gridInfoMap` field does not populate member list reliably; raw `memberIds` is often empty. | Use **`api.getGroupMembersInfo(groupId)`** instead — but it too may return partial data. **Fallback:** call `getAllGroups()` → iterate `gridVerMap` keys (the stable group ID list), then `getGroupInfo(id)` per group. For member detection across sessions, maintain your own `group_members_cache.json`. |
| **Cookie login fails after forceful kill** | Killing the node server process without calling `client.shutdown()` leaves the Zalo socket in a half-open state on Zalo's side → next login with same cookie returns `Đăng nhập thất bại`. | Before killing/restarting the bridge: (1) hit `POST /shutdown` endpoint on the bridge, OR (2) run `node -e "const{Zalo}=require('zca-js');new Zalo().login(require('./.hermes-zalo/credentials.json')).then(c=>c.shutdown())"` — this gracefully tears down the socket so the next cookie login succeeds. |
| **Bridge shows `loggedIn: true` but SSE not receiving messages** | Gateway SSE client disconnected because `sseClients: 0` in `/health` — usually because gateway restarted before bridge finished QR re-login (race condition). | Always wait ≥8 s after `pm2 resurrect` before checking `/events` or expecting gateway reconnect. Restart gateway (`hermes gateway restart`) **after** bridge is `loggedIn: true`. |

## ⚠️ Known Limitations

- **Zalo kick on duplicate connection**: Running two bridge instances with the same account will kick the previous one. Use PM2's `resurrect` or `monitor` to auto-restart after crash, but avoid running two simultaneously.
- **GoClaw / OX Gateway Bridge Pattern**: Instead of embedding `zca-js` or unofficial Zalo protocol into every gateway or bot agent (which causes session conflict code 3000 / kicks), connect agents (GoClaw, OX Alpha, etc.) to the local bridge HTTP port (e.g. `http://127.0.0.1:8787`):
  - **Inbound**: Connect to SSE stream `GET /events` to receive incoming user/group messages.
  - **Outbound**: Send `POST /send` with `{ threadId, threadType, text }`.
  - **File/Media**: Send `POST /send-attachment` with `{ threadId, threadType, filePath }`.
- **GoClaw Direct Native Zalo Implementation Reference**: If embedding Zalo directly into a Go binary without external Node bridge is needed, refer to GoClaw's implementation (`github.com/nextlevelbuilder/goclaw`):
  - **Zalo OA**: `internal/channels/zalo/` uses official Zalo Bot API (`https://bot-api.zaloplatforms.com`) via webhook/polling (DM only).
  - **Zalo Personal**: `internal/channels/zalo/personal/protocol/` ports `zcago` protocol directly to Go, connecting via `gorilla/websocket` with manual cookie jar injection to `wss://` endpoints and AES-CBC/RSA crypto routines.
- **Management vs Gateway MCP Distinction**: Management MCP servers (e.g. `goclaw-mcp`) do NOT connect to chat channels (Zalo/Discord); they provide external IDE clients (Claude Code/Cursor) with tools to control and configure the gateway via REST API endpoints.
- **Cookie expiration**: Zalo session cookies typically live 1–7 days. After expiry, the bridge will fall back to QR re-scanning.
- **No native `welcome` event**: The polling-based welcome (15s interval) is a workaround; frequency can be adjusted in `server.js`.
- **Member detection via `getGroupInfo`**: Relies on `gridInfoMap.memberIds` which Zalo may occasionally reset or truncate.

---

## 🔧 Troubleshooting

| Issue | Fix |
|---|---|
| **Bot not responding in Zalo** | Ensure bridge is running (`pm2 status`), gateway has `allowed_threads` / `allowed_groups` set, and `hermes gateway restart` was executed. |
| **Connection refused on port 8787** | Kill any stale node processes: `powershell.exe -NoProfile -Command "Get-Process node \| Stop-Process -Force"`. Restart bridge. |
| **Members not being tracked** | Ensure the poller is running (`pm2 logs monitor`) and the target group ID matches `allowed_threads` config. |
| **Bot replies to itself** | Check `uidFrom` vs bot's own UID (`37082029033239591`) and filter in logic. |
| **Bot fails to send/upload files** | 1. Hardcoded `ALLOWED_THREADS` in `server.js` (line ~350) filters text replies (`/send`), but verify if endpoint used is `/send-attachment`.<br>2. Local file paths must exist directly on the machine running the bridge process (absolute or accessible path).<br>3. Bridge policy must include `uploadAttachment` in `ZALO_ALLOWED_ACTIONS` or `send` in `ZALO_ALLOWED_ACTION_GROUPS`. |

---

## 📦 Package References

- `hermes-zalo-plugin` (v1.0.9) — external Node.js bridge
- `zalo-platform` (Hermes plugin) — Hermes gateway adapter
- `zca-js` — unofficial Zalo API client used by the bridge
- `pm2` — process manager for auto-restart

## 🧰 scripts/get_zalo_group_ids.mjs
**Purpose**: Discover Zalo group IDs using in-process `zca-js` login (cookie or QR fallback).

```javascript
// scripts/get_zalo_group_ids.mjs
import { ZaloClient } from "../zaloClient.js";
import path from "node:path";

async function getGroupIds() {
  const c = new ZaloClient({
    credentialsPath: path.join(process.env.HOME, ".hermes-zalo", "credentials.json"),
    qrPath: path.join(process.env.HOME, ".hermes-zalo", "qr.png"),
  });
  await c.login();
  const result = await c.api.getAllGroups();
  const groupIds = Object.keys(result?.gridVerMap || {});
  console.log("=== ZALO GROUP IDs ===");
  for (const id of groupIds) {
    try {
      const info = await c.api.getGroupInfo(id);
      const name = info?.gridInfoMap?.[id]?.name || "N/A";
      console.log(`${name}\t| ${id}`);
    } catch (e) {
      console.log(`?\t| ${id}`);
    }
  }
  process.exit(0);
}

getGroupIds().catch((e) => {
  console.error("Error:", e.message);
  process.exit(1);
});
```

Run with:
```bash
cd C:/Users/OS/AppData/Local/hermes/plugins/hermes-zalo-plugin
node scripts/get_zalo_group_ids.mjs
```