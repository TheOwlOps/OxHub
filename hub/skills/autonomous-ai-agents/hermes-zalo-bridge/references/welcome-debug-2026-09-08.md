# Zalo welcome debug — 2026-09-08 (group 8813682027038154228)

Symptom: member join/leave visible in Zalo client, no bot welcome in group.

## Root cause 1 (blocking): adapter.py IndentationError disables plugin
- `adapter.py` had class-level code at indent 0 with `await self.handle_message(event)`
  outside any `async def` → `SyntaxError: 'await' outside function (line 647)`.
- Effect: `discover_plugins(force=True)` records
  `hermes-zalo-plugin/hermes-plugin enabled=False`, gateway `load_gateway_config()`
  shows only `telegram`, bridge `/health` shows `sseClients: 0`.
- Detect fast (no gateway restart needed):
  `ZALO_PLUGIN_URL=http://127.0.0.1:8787` then `discover_plugins(force=True)` and
  inspect `pm._plugins['hermes-zalo-plugin/hermes-plugin'].error`.
- Fix: keep all inbound logic inside `async def _on_inbound_message`,
  `await self.handle_message(event)` indented at 8 spaces, `async def _download_media`
  at 4 spaces. Verify with `python3 -m py_compile adapter.py` before `pm2 restart monitor`.

## Root cause 2 (silent): getGroupInfo memberIds empty
- `GET /chat-info?threadId=8813682027038154228&threadType=group` and
  `POST /api/getGroupInfo` both returned `memberIds: []`, `currentMems: []`,
  `adminIds: []` while `memVerList: ["37082029033239591_0","6171736846235386882_0"]`
  and `totalMember: 2` were correct.
- Poller in `server.js:checkNewMembers` diffed only `memberIds` → cache file
  `~/.hermes-zalo/group_members_cache.json` stayed `{"8813682027038154228": []}`,
  no `welcome-poller` lines in `~/.pm2/logs/monitor-out.log`.
- Fix: parse members from `memVerList` (`uid.split('_')[0]`) or fall back to
  `totalMember` delta, never `memberIds` alone.

## Root cause 3 (silent): `threadId` read from wrong field (2026-09-08 PM)

- `adapter.py` did `thread_id = data.get("threadId")` on group_event — but
  group_event payload has NO top-level `threadId`. Verified event shape from
  live `gateway.log`:
  ```json
  {"type": "join", "data": {"groupId": "8813682027038154228",
   "groupName": "Hermes ( Hơi Mệt)",
   "updateMembers": [{"id": "7378254928662738368", "dName": "Sự"}]}}
  ```
- Symptom: `Zalo: raw group_event data: {... type join ...}` logged, but no
  `Zalo: chào mừng thành viên mới` line and no message sent — `thread_id`
  was `""` so the send block was skipped silently.
- Fix: `thread_id = str(data.get("threadId") or group_data.get("groupId") or "")`.

## Group name: read from event, never hardcode (2026-09-08 PM)

- Welcome text hardcoded `"Hermes (Hơi Mệt)"` — group rename would not follow.
  Fix: `group_name = group_data.get("groupName") or "nhóm này"`, interpolate
  into the welcome text. Verified `groupName` present on every join event.

## Multi-UID joins (2026-09-08 PM)

- Old code only greeted `new_members[0]`. Zalo can add several users in one
  event (verified: 2 UIDs in one join payload). New code greets all UIDs in
  one message: text prefix `"@member " * n`, mentions at pos 0,8,16...
  (`@member` = 7 ASCII chars + 1 space). Log line now records the UID list:
  `Zalo: chào mừng thành viên mới [uids] vào nhóm <id> (<group_name>)`.

## Gateway restart pitfalls (2026-09-08 PM)

- `hermes gateway restart` is BLOCKED when run inside the gateway process
  (SIGTERM kills the caller). Run from an external shell instead, e.g.
  `powershell.exe -NoProfile -Command "Start-ScheduledTask -TaskName Hermes_Gateway"`.
- `adapter.py` loads once at gateway start — code edits need a restart to go live.
- After restart, `sseClients` may read 0 for a few seconds (bridge still logging
  in). Wait ≥8 s before concluding failure.
- Multiple stale `python.exe gateway run` processes cause `sseClients: 0`
  while port stays open. Kill stale PIDs, then `Start-ScheduledTask`.
- Auto-restart watcher (PM2 `zalo-watcher`, CommonJS file must be `.cjs`
  because plugin `package.json` has `"type": "module"`): `fs.watchFile`
  on `adapter.py` (poll 2 s, debounce 15 s) → `Start-ScheduledTask`.
  Location: `hermes-zalo-plugin/watch_zalo_adapter.cjs`, log `watcher.log`.
- Plain text `@All` renders black. Blue ping-all needs mention payload through
  `POST /send`: `{"threadId","threadType":"group","text":"@All ...",
  "mentions":[{"pos":0,"uid":"-1","len":4}]}`.
- Verified `msgId: 8237920456401` rendered blue. `zaloClient.js:handleMentions`
  maps `uid "-1"` → `type:1`.
- Welcome mention for one user: `text "Chào mừng @member ..."`,
  `mentions:[{"pos":10,"uid":uid,"len":7}]`.

## Checks
- `curl -s http://127.0.0.1:8787/health` → want `loggedIn:true`, `sseClients:1` when
  gateway adapter attached (`0` = adapter down, not bridge bug).
- `pm2 logs monitor --lines 30 --nostream` → look for `[welcome-poller]`,
  `[zalo] RAW message`, `[zalo] group_event`.
- Cache: `~/.hermes-zalo/group_members_cache.json` (via `paths.js:dataDir()`).
- Gateway log `gateway.log` showed only `telegram` (`Gateway running with 1 platform(s)`)
  while Zalo was broken — expected, not a separate bug.
