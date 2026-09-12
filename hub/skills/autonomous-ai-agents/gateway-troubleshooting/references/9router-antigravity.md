# 9Router Antigravity OAuth Refresh

Local proxy `9Router` (`C:/Users/OS/AppData/Roaming/9router/db.json`, also `$APPDATA/9router/db.json`) stores OAuth providers in `providerConnections[]` where `provider == "antigravity"`.

## Inject a new refresh token

1. Fix prefix — 9Router expects `1//...`. Raw paste from Google OAuth often lacks it:
   ```python
   raw = "0eZZ..."  # pasted from user
   fixed = raw if raw.startswith("1//") else "1//" + raw
   ```
2. Edit db.json atomically:
   ```python
   import json
   p = r"C:/Users/OS/AppData/Roaming/9router/db.json"
   d = json.load(open(p, encoding="utf-8"))
   for c in d["providerConnections"]:
       if c["provider"] == "antigravity":
           c["refreshToken"] = fixed
           c.pop("accessToken", None)
           c["expiresAt"] = "2020-01-01T00:00:00.000Z"  # force refresh on next tick
           c["isActive"] = True
           c["backoffLevel"] = 0
           c["lastError"] = c["errorCode"] = c["testStatus"] = None
           for k in list(c):
               if k.startswith("modelLock_"): del c[k]
   json.dump(d, open(p, "w", encoding="utf-8"), indent=2, ensure_ascii=False)
   ```
3. Restart the proxy — RAM caches `backoffLevel` and `modelLock_*` until restart:
   ```bash
   taskkill /PID <pid> /F   # pid from: netstat -ano | findstr 20128
   # or: wmic process where "ProcessId=<pid>" get CommandLine
   # restart via the same launch that was used (node custom-server.js --port 20128)
   ```
   Verify: `curl -s http://localhost:20128/v1/models -H "Authorization: Bearer <sk>" | head -c 500` and `curl .../v1/chat/completions -d '{"model":"ag/claude-sonnet-4-6",...}'`

## Failure modes

- `429 RESOURCE_EXHAUSTED / Individual quota reached / modelLock_*` — account quota exhausted. Reset timer is in the error metadata (`quotaResetDelay`, e.g. `115h...`). Injecting a new refresh token for the same Google account does not help if the quota is per-account; add a different Google account or wait for reset.
- `403 VALIDATION_REQUIRED` — Google anti-abuse checkpoint. Prompts SMS to US shortcode `96831`. Carriers outside US fail to send. Account dead for Antigravity until checkpoint lifts; disable connection in DB/UI.
- `403 Organization Restriction` — Domain/edu account blocked by workspace admin (`"restricted from using Gemini Code Assist for individuals in your organization"`). Must switch to personal `@gmail.com`.
- `404 Requested entity was not found` — Model name mismatch in test ping or inactive project ID (`aicode-consumers` vs `tangential-setup-*`).
- `405 Method Not Allowed` with `reset after 30s` — upstream rate-limit window, retry after window.
- `invalid_client` on `https://oauth2.googleapis.com/token` — wrong `client_id` for that refresh token (each provider has its own OAuth client, not the generic Google one).
