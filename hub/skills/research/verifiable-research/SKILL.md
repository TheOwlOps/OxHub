---
name: verifiable-research
description: "Fetch live data via scrape/API with citations."
version: 1.0.0
author: Hermes Agent + curator
license: MIT
platforms: [linux, macos, windows]
metadata:
  hermes:
    tags: [research, verification, trending, scraping, time, grounding]
    category: research
    related_skills: [grounded-citations, competitor-news-monitor, arxiv]
---

# Verifiable Research

Every live-data answer is fetched at call time, parsed from tool output, and cited. No memory, no guessing, no model-knowledge fill-ins.

## When to Use

- User asks for "trending today/this week", "top repos", "price now", "what time is it"
- Any claim that changes daily (market, rankings, news, rates)
- User has corrected a hallucination before (`bịa kết quả`, `sai ngày`, `sai giờ`)
- MSYS2/Windows toolchain blocks a fancier tool — fall back to stdlib fetch

Skip for: pure coding questions, static docs, creative writing.

## Prerequisites

- `python3` + stdlib `urllib.request`, `re`, `json`, `datetime` (no extra pip needed)
- Optional: `hermes-supersearch` / `browser-act` only if already installed; never block on install
- For MCP browser: `chrome-devtools` enabled via `hermes mcp add` — requires gateway restart

## Procedure

### 0. Reset & choose source (before any prose)

Ledger optional — if grounded-citations available, register URLs there; otherwise just fetch directly; never invent a URL.

### 1. Time — always ground first

```bash
date '+%A, %d/%m/%Y %H:%M:%S (GMT+7)'
# fallback when date lies (WSL clock drift):
python3 -c "import datetime; print(datetime.datetime.now().astimezone().strftime('%A, %d/%m/%Y %H:%M:%S (%Z)'))"
```

Rules:
- Run the command on every time question, even if asked twice in a row.
- Never answer from context/memory. Memory rule `BẮT BUỘC: date trước khi trả lời` is the ceiling — this skill implements it.
- If `date` and `Get-Date` disagree, trust the newer `date` after `tzutil /s "SE Asia Standard Time" && w32tm /resync` (requires Admin).

### 2. GitHub Trending — scrape + enrich (works on MSYS2 without extra deps)

```bash
python3 << 'PY'
import urllib.request, re, json
url = "https://github.com/trending?since=daily"  # or weekly/monthly
headers = {"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"}
raw = urllib.request.urlopen(urllib.request.Request(url, headers=headers), timeout=20).read().decode("utf-8", errors="ignore")
pat = re.compile(r'<h2[^>]*class="[^"]*h3[^"]*"[^>]*>.*?<a[^>]*href="(/[^/]+/[^/"]+)".*?>(.*?)</a>', re.DOTALL)
repos = []
for href, name_raw in pat.findall(raw):
    href = href.split("?")[0]
    if href.count("/") != 2 or "/trending" in href: continue
    name = re.sub(r'<[^>]+>', '', name_raw).strip().replace(" ","").replace("\n","")
    name = re.sub(r'\s*/\s*', '/', name)
    if "/" not in name: name = href.lstrip("/")
    repos.append(name)
    if len(repos) >= 10: break
for r in repos[:10]:
    data = json.loads(urllib.request.urlopen(urllib.request.Request(
        f"https://api.github.com/repos/{r}",
        headers={"User-Agent":"Mozilla/5.0","Accept":"application/vnd.github.v3+json"})).read().decode())
    print(f"{r} | {data['stargazers_count']:,}⭐ | {data['language']} | {data['description'][:90]}")
PY
```

- Scrape gives order (trending rank); API gives absolute stars. Quote both.
- `since=weekly` → `?since=weekly` for weekly top 10; `since=monthly` similarly.

### 3. Gold / market prices — multi-source scrape + compare

```bash
python3 -c "
import urllib.request, re
url='https://24h.com.vn/gia-vang-hom-nay-c425.html'
html=urllib.request.urlopen(urllib.request.Request(url, headers={'User-Agent':'Mozilla/5.0'})).read().decode('utf-8', errors='ignore')
rows=re.findall(r'<tr[^>]*>(.*?)</tr>', html, re.DOTALL)
for row in rows[:15]:
    cells=[re.sub(r'<[^>]+>','',c).strip() for c in re.findall(r'<td[^>]*>(.*?)</td>', row, re.DOTALL)]
    if any(cells): print(' | '.join(cells))
"
```

- Primary: `24h.com.vn` (real-time per brand SJC/DOJI/BTMC). Secondary: `giavang.net`, `pnj.com.vn`, `sjc.com.vn` (often 403 — note as unavailable, don't fake).
- When sources disagree heavily, state both and flag the likely stale one; never smooth over.
- Always include fetch timestamp and source URL in answer.

### 4. If a nicer tool exists, prefer stdlib fallback

- `hermes-supersearch` on MSYS2 fails with `externally-managed-environment` + `cmake/ninja` + `numpy` wheel build. Fix is venv, not global pip:
  ```bash
  python3 -m venv ~/venvs/hermes-websearch
  source ~/venvs/hermes-websearch/bin/activate
  pip install --only-binary :all: git+https://github.com/hermes-labs-ai/supersearch.git
  ```
- `browser-act-cli` requires `cp312` wheels — fails on `cp311` MSYS Python. Requires `uv tool install browser-act-cli --python 3.12`.
- `chrome-devtools-mcp` needs `hermes mcp add chrome-devtools --command npx --args -y chrome-devtools-mcp@latest` then gateway restart **outside** gateway (`hermes gateway restart` from separate shell — inside-gateway restart is blocked).
- Until those are fixed, the urllib scrape above is the reliable path. Don't present install failures as "tool broken" — present the venv fix.

### 5. Browser vs API vs computer_use — choose by least chrome tabs

- Least disruptive first: `urllib` scrape/API (no Chrome tab).
- Next: `chrome-devtools` MCP `navigate_page` + `take_snapshot` (headless, 29 tools) — only if page needs JS rendering.
- Last: `computer_use` capture/click — foreground delivery required for `Chrome_WidgetWin_1` (background_unavailable). Avoid if user complained about tab popups.

### 6. Cite while drafting

- After each fetch, cite the URL. Max 3 ids per sentence. No space before `[n]`.
- For chat, end with `Sources:` list rendered from the fetches you just did.
- When listing GitHub trending, render each repo as a tappable markdown link `[owner/repo](https://github.com/owner/repo)` — user explicitly corrected plain-text lists to require clickable links.
- If a value truly not found, say `no source found for X` — never fill with model knowledge.

## Pitfalls

- **Guessing time from context** — 14-month drift observed (Saturday 05/09/2026 vs Friday 04/09/2026). Fix: always run `date`; never reuse last answer.
- **Hallucinating 500M VND cars** — listed Genesis GV80/BMW X3/GLC (2-4B VND) as 500M. Fix: scrape real price tables, hallucinate = immediate trust loss.
- **Single-source gold** — `giavang.net` returned 249M vs 24h's 143M due to stale markup. Fix: fetch 2 sources, compare, flag outlier.
- **Trending via Search API** — `api.github.com/search/repositories?q=created:>date` misses trending (stars velocity not creation date). Fix: scrape `/trending` HTML + enrich via API.
- **Trending via web_search** — aggregator caches (findarepo, wangchujiang, unpkg) go stale and mismatch the live page (observed 08/09 cache vs 09/09 live). Fix: scrape `github.com/trending` HTML at call time, never quote a search snippet as the ranking.
- **Trending list without links** — plain `owner/repo` text is not tappable. Fix: always emit `[owner/repo](https://github.com/owner/repo)`.
- **MCP session stale** — `navigate_page pageId 1 → No page found` after restart. Fix: `list_pages` then use returned id; or use `new_page` instead of assuming id 1.
- **Inside-gateway restart blocked** — `hermes gateway restart` from inside returns `Blocked: ... would kill this command`. Fix: run from separate PowerShell/shell outside gateway, or Task Scheduler `Hermes_Gateway → Run`.
- **MSYS2 pip externally-managed** — `pip install --user` blocked by PEP 668. Fix: `python3 -m venv ~/venvs/<name>` + activate, or `uv pip install --break-system-packages` on that venv only.
- **Retyping URLs** — always copy from fetch headers, never retype.
- **GitHub repo detail via web_extract empty** — github.com pages return a JS shell with no extractable content. Fix: fetch `https://raw.githubusercontent.com/{owner}/{repo}/main/README.md` via curl and read the first ~100 lines for purpose, install, and usage.
- **TikTok trending via web_extract thin** — /trending and /discover are client-rendered, so extract returns a shell or stale tags. Fix: snapshot the page in a real browser and collect `a[href*="/tag/"]` innerText for live hashtags.
- **Hallucinating YouTube/media video IDs** — LLMs fabricate plausible 11-char strings (e.g. `watch?v=UpJyBAh8JS4` instead of real `Vt4kAu-ziRY`) with >50% failure rate. Fix: scrape `https://www.youtube.com/results?search_query=<encoded>` and extract real `"videoId":"([a-zA-Z0-9_-]{11})"` regex matches at call time; if scraping fails, emit the search URL (`https://www.youtube.com/results?search_query=...`) rather than guessing a broken video link.

## Verification

```bash
date '+%A, %d/%m/%Y %H:%M:%S (GMT+7)'
python3 -c "import urllib.request, re; raw=urllib.request.urlopen(urllib.request.Request('https://github.com/trending', headers={'User-Agent':'Mozilla/5.0'})).read().decode(); print(len(re.findall(r'/[^/]+/[^/]+', raw)))"
curl -I https://24h.com.vn/gia-vang-hom-nay-c425.html
```

## References

- `references/github-trending.md` — HTML structure + API enrich recipe that survived MSYS2 constraints
- `references/gold-scraping.md` — 24h.com.vn regex + multi-source comparison template
