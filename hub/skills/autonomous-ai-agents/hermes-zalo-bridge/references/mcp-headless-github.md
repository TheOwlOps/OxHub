# MCP headless + GitHub — session 2026-09-08

## Chrome-devtools headless (user hates visible tab)

Config path: `C:/Users/OS/AppData/Local/hermes/config.yaml`

Desired block:
```yaml
mcp_servers:
  chrome-devtools:
    command: npx
    args:
      - -y
      - chrome-devtools-mcp@latest
      - --headless
      - --isolated
    enabled: true
```

Fix steps when `list_pages` fails with "browser already running":
1. Kill stale Chrome: `taskkill /F /IM chrome.exe`
2. Delete profile lock: `rmdir /S /Q "%USERPROFILE%\.cache\chrome-devtools-mcp\chrome-profile"`
3. Restart gateway OUTSIDE gateway: `powershell Start-ScheduledTask -TaskName Hermes_Gateway` (never `hermes gateway restart` inside gateway — blocked, SIGTERM kills caller)
4. Verify: `hermes gateway status` + `list_pages` should return new page id

Pitfall: `hermes mcp rm chrome-devtools` needs explicit user confirm. User said "t bảo m chạy nền chứ đâu kêu m xoá" — never delete without ask. Prefer editing args to add headless.

Research fallback when MCP blocked:
- Google `ERR_BLOCKED` / captcha → use Bing: `https://www.bing.com/search?q=...` works via MCP
- `giavaxang.vn` → DNS fail, `price.24h.com.vn` → timeout → Bing search is stable
- Alternative: `curl` / python `urllib` for API scrape when browser not needed

## GitHub MCP server (26 tools)

Token: `ghp_...` for user `ryanowlops` (8 repos). Works via stdio:
```bash
echo "Y" | hermes mcp add github --command "npx" --args "-y" "@modelcontextprotocol/server-github" --env "GITHUB_PERSONAL_ACCESS_TOKEN=<PAT>"
# verify
hermes mcp list
powershell Start-ScheduledTask -TaskName Hermes_Gateway  # restart outside
```
Tools: list_issues, search_issues, create_issue, get_issue, create_pull_request, search_code, list_commits, etc. Called as `mcp__github__<tool>`.

Repo delete via API (used for vibe-bypass, BunVibe):
```bash
curl -s -X DELETE -H "Authorization: Bearer $PAT" https://api.github.com/repos/ryanowlops/<repo>
```
List: `GET /users/ryanowlops/repos?sort=updated&per_page=30`

## Gateway restart rule

Inside gateway process, `hermes gateway restart` is blocked. Always use `powershell Start-ScheduledTask -TaskName Hermes_Gateway` + 8-10s wait + `hermes gateway status`.
