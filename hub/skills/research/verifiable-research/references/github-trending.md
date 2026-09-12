# GitHub Trending — scrape recipe (MSYS2-safe)

Validated 2026-09-07 on Windows MINGW64, Python 3.14, no extra pip.

## Why not Search API

`GET /search/repositories?q=created:>2026-09-04&sort=stars` finds recently *created* repos, not trending (velocity). Trending rank comes only from `https://github.com/trending` HTML.

## HTML structure (stable Aug 2026)

```html
<article class="Box-row">
  <h2 class="h3 lh-condensed">
    <a href="/owner/repo">owner / repo</a>
  </h2>
  ...
</article>
```

Regex that survived:
```python
pat = re.compile(r'<h2[^>]*class="[^"]*h3[^"]*"[^>]*>.*?<a[^>]*href="(/[^/]+/[^/"]+)".*?>(.*?)</a>', re.DOTALL)
```

Filters:
- `href.count("/") == 2`, skip `/trending`, `/login`, `/features`
- Clean name: strip tags, remove spaces/newlines, `re.sub(r'\s*/\s*','/',name)`, fallback to `href.lstrip("/")`

Daily: `https://github.com/trending`
Weekly: `https://github.com/trending?since=weekly`
Monthly: `https://github.com/trending?since=monthly`

## Enrich via API (true stars/lang/desc)

```python
import urllib.request, json
data = json.loads(urllib.request.urlopen(urllib.request.Request(
    f"https://api.github.com/repos/{repo}",
    headers={"User-Agent":"Mozilla/5.0","Accept":"application/vnd.github.v3+json"})).read().decode())
stars, lang, desc = data["stargazers_count"], data["language"], data["description"]
```

Rate limit: unauthenticated 60/h. Trending = 10 calls, safe. Auth with `Authorization: Bearer <PAT>` if needed.

## Verified examples

- 2026-09-07 weekly trending top 3: `magnitudedev/magnitude | 3,981⭐`, `tt-a1i/archify | 52,291⭐`, `Gitlawb/openclaude | 32,908⭐`
- 2026-09-07 daily top 3: `heygen-com/hyperframes`, `microsoft/markitdown`, `mksglu/context-mode`

## Failure seen

- SciPy-family stacks tried via `pip --user` on MSYS2 hit `externally-managed-environment` (PEP 668) + `cmake`/`ninja` wheel build failures. Don't block trending on that — use stdlib recipe above.
