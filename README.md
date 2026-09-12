# OxHub ⚡

> **Universal Agent Hub & Ecosystem** — High-performance statusline HUD, skills registry, and plugin ecosystem for **Claude Code**, **Hermes Agent**, and **OpenAI Codex**.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://golang.org)
[![Cross-Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-brightgreen)](#)

---

## ⚡ 1-Line Universal Install (Community & Agents)

Install all 120+ skills and tools directly into your agent environment with a single command:

### Linux / macOS:
```bash
curl -fsSL https://raw.githubusercontent.com/TheOwlOps/OxHub/main/install.sh | bash -s -- all
```
*(Or target a specific agent: `bash -s -- hermes`, `bash -s -- claude`, `bash -s -- codex`)*

### Windows (PowerShell):
```powershell
irm https://raw.githubusercontent.com/TheOwlOps/OxHub/main/install.ps1 | iex
```

---

## 🎯 Native CLI Commands (Per-Agent Setup)

### 1. Hermes Agent
Hermes natively manages external skills via `skills tap` and `skills install`:

- **Add the full OxHub repository as an official tap**:
  ```bash
  hermes skills tap add TheOwlOps/OxHub
  ```
- **Install individual skills via ID or direct URL**:
  ```bash
  hermes skills install TheOwlOps/OxHub/hub/security/security-audit
  # Or via direct RAW URL:
  hermes skills install https://raw.githubusercontent.com/TheOwlOps/OxHub/main/hub/security/security-audit/SKILL.md
  ```

---

### 2. Claude Code
- **Install via Plugin Marketplace**:
  ```bash
  /plugin marketplace add TheOwlOps/OxHub
  /plugin install oxhud
  ```
- **Automated StatusLine Configuration**:
  Add `statusLine` to your `~/.claude/settings.json`:
  ```json
  {
    "statusLine": "oxhud --agent=claude"
  }
  ```
- **Automatic Skills Discovery**:
  All skills copied into `~/.claude/skills/` are automatically discovered and indexed by Claude Code in every session.

---

### 3. OpenAI Codex
- Run the wrapper or load skills directly from `~/.codex/skills/`:
  ```bash
  echo '{"model":"gpt-4o","used_tokens":110000,"total_tokens":128000,"active_tool":"edit_file","status":"patching"}' | oxhud --agent=codex
  ```

---

## 🚀 OxHUD (CLI Real-time Statusline)

A zero-dependency, ultra-low latency statusline written in Go:

```text
[CLAUDE: Opus 3.7] ██░░░░░░░░ 50k/200k (25.0%)
◐ Tool: edit_file | Task: patching state.go
```

### Build from Source:
```bash
go build -o bin/oxhud ./cmd/oxhud
```

---

## 📦 Skills & Plugins Catalog (120+ Skills)

All skills are declared and structured in `hub/manifest.json`.

| Category | Featured Skills | Highlights |
| :--- | :--- | :--- |
| 🛡️ **Security** | `security-audit`, `pen-test`, `security-and-hardening` | Secret scanning, SQLi/RCE/SSRF mitigation, OWASP Top 10 & API pentest checklists. |
| 💻 **Coding** | `code-reviewer`, `systematic-debug`, `refactor-cleanup` | Multi-axis code reviews, YAGNI enforcement, 4-phase root-cause debugging. |
| 🤖 **Autonomous** | `claude-code`, `hermes-agent`, `codex`, `computer-use` | Multi-agent delegation, worktree isolation, desktop & browser automation. |
| 🧠 **MLOps** | `evaluating-llms-harness`, `llama-cpp`, `serving-llms-vllm`, `weights-and-biases` | GGUF quantization, vLLM deployment, evaluation harnesses, W&B sweeps. |
| 📊 **Productivity**| `google-workspace`, `pdf`, `docx`, `xlsx`, `airtable`, `notion` | Document processing pipelines, spreadsheet automation, headless office generation. |

---

## 📄 License
Distributed under the [MIT License](LICENSE).
