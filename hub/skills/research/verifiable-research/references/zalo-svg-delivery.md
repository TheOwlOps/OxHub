# Zalo SVG Delivery — HTML vs SVG

Zalo bridge (`POST /send-attachment`) renders `*.svg` inline as a previewable image; `*.html` is sent as a downloadable file with no preview (user must tap to download).

## Rule
- When source is an infographic/diagram HTML (e.g. `chinese_sach_vs_doi.html`), convert to standalone SVG (`width=1200`, inline styles, no external CSS) before sending to Zalo.
- Keep `MEDIA:` reply for Telegram (supports both), but Zalo path must be `*.svg`.

## Pitfall
- Sending HTML to Zalo group loses preview and forces download. Fix: always send `*.svg` via `send-attachment` for diagrams; reserve HTML for local preview only.

Validated 2026-09-09: `chinese_sach_vs_doi.svg` (5.6KB) sent via curl to `8813682027038154228` succeeded with `msgId`, `pnj.com.vn`-style HTML failed preview.
