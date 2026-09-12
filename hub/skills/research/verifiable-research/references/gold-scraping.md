# Gold scraping — 24h.com.vn + comparison

Validated 2026-09-07 23:26 GMT+7.

## Primary source

`https://24h.com.vn/gia-vang-hom-nay-c425.html` — real-time per brand table. Works with `urllib` + desktop UA; no JS required.

```python
import urllib.request, re
url='https://24h.com.vn/gia-vang-hom-nay-c425.html'
html=urllib.request.urlopen(urllib.request.Request(url, headers={'User-Agent':'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'})).read().decode('utf-8', errors='ignore')
rows=re.findall(r'<tr[^>]*>(.*?)</tr>', html, re.DOTALL)
for row in rows:
    cells=[re.sub(r'<[^>]+>','',c).strip() for c in re.findall(r'<td[^>]*>(.*?)</td>', row, re.DOTALL)]
    if any(cells): print(' | '.join(cells))
```

Sample output (07/09/2026):
```
SJC | 143,500 1100 | 146,500 1100 | 144,600 | 147,600
DOJI HN | 143,500 1100 | 146,500 1100 | 144,600 | 147,600
BTMH | 144,600 900 | 148,600 900 | 145,500 | 149,500
```

## Secondary sources to compare

- `https://giavang.net/` — often stale; on 2026-09-07 returned `249,259,242 / 252,219,229` for SJC vs 24h's `143,500 / 146,500` (flag as outlier, don't smooth).
- `https://www.sjc.com.vn/giavang/textContent.php` — returned 403 Forbidden (note unavailable, don't fake).
- `pnj.com.vn`, `doji.vn` — candidates for future comparison.

## Comparison template

| source | SJC buy | SJC sell | note |
|--------|---------|----------|------|
| 24h.com.vn | 143,500 | 146,500 | real-time, cited |
| giavang.net | 249,259 | 252,219 | stale — flag outlier |
| sjc.com.vn | — | — | 403 unavailable |

## Time-grounding tie-in

Gold prices are intraday — always include fetch time `date '+%A, %d/%m/%Y %H:%M:%S (GMT+7)'`. User corrected 14-month drift (05/09/2026 Saturday vs actual Monday 07/09/2026 after tzutil fix), so time citation is load-bearing.

## Tooling note

No npm/pip required. Stdlib regex survives MSYS2 `externally-managed-environment`. Fancier `hermes-supersearch` failed with numpy/cmake/ninja build — fallback is this recipe.
