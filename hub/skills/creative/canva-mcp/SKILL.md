---
name: canva-mcp
category: creative
description: "Work with Canva via MCP; always export to media files."
version: 1.0.0
author: Hermes Agent
license: MIT
platforms: [linux, macos, windows]
metadata:
  hermes:
    tags: [Canva, MCP, Design, Export]
    related_skills: [media-broadcast]
---

# Canva MCP Workflow

## When to use

Use when users request: Canva design creation, logo generation, design editing, template use, design export, or sharing Canva designs.

## Core Rule (NON-NEGOTIABLE)

**NEVER send raw Canva share URLs (design.canva.ai/..., canva.com/d/..., export-download.canva.com/...).** They are:
- Private by default (403 Forbidden)
- S3-signed with short expiry (hours)
- Contain credentials in URL (leaks if shared)
- Truncated/escaped by Discord/Telegram

**ALWAYS:**
1. `mcp__canva__export_design` with `format: {type: "png"}` (or pdf)
2. Download the returned S3 URL to local file (`curl -o file.png "$URL"`)
3. Send as `MEDIA:/absolute/path/to/file.png` in response

## Procedure

### 1. Generate design candidates
```
mcp__canva__generate_design({
  design_type: "logo|presentation|social-media|...",
  query: "detailed prompt",
  user_intent: "what user wants"
})
```
Returns job_id + candidate_ids (dg-...).

### 2. Convert candidate → real design
```
mcp__canva__create_design_from_candidate({
  candidate_id: "dg-...",
  job_id: "..."
})
```
Returns design_summary with `id` (DAHU...), `view_url`, `edit_url`.

### 3. Export to file
```
mcp__canva__export_design({
  design_id: "DAHU...",
  format: {type: "png"}
})
```
Returns job with `urls[0]` = S3 signed URL.

### 4. Download & send
```bash
curl -sL "$S3_URL" -o "C:/Users/OS/design-name.png"
```
Then in response: `MEDIA:C:/Users/OS/design-name.png`

## Pitfalls

- **Raw URL sharing fails**: Canva view/edit URLs require login; export URLs expire and contain AWS creds. User sees 403/blank page.
- **Discord truncates long URLs**: S3 signed URLs exceed Discord char limit → broken link.
- **Zalo bridge rejects HTML**: User configured Zalo to only accept image files, not HTML/download links.
- **Candidate ≠ Design**: Candidates are preview only. Must `create_design_from_candidate` first to get persistent design ID.
- **MCP auth**: Requires Canva OAuth connected in Hermes config. If `get_design` fails, check MCP connection status.

## Decision Table

| Goal | MCP Call | Output | Next Step |
|------|----------|--------|-----------|
| New design idea | `generate_design` | candidates (dg-...) | Pick one → `create_design_from_candidate` |
| Use existing design | `get_design` / `search_designs` | design_summary (DAHU...) | `export_design` |
| Share with others | `export_design` (png/pdf) | S3 URL | `curl` download → `MEDIA:` |
| Edit design | `start_editing_transaction` + `perform_editing_operations` + `commit` | Updated design | `export_design` |
| Copy template | `copy_design` | New design ID | Edit or export |

## Verification

Before sending to user:
- [ ] File exists locally (`ls -lh file.png`)
- [ ] File is valid image (`file file.png` → PNG image data)
- [ ] Sent as `MEDIA:` not markdown link
- [ ] No Canva URLs in response text