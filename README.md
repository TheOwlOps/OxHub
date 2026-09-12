# OxHub ⚡

> **Universal Agent Hub & Ecosystem** — A high-performance, open-source statusline HUD, skills registry, and plugin ecosystem for **Claude Code**, **OpenAI Codex**, and **Hermes Agent**.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://golang.org)
[![Cross-Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-brightgreen)](#)

---

## 🌟 What is OxHub?

OxHub delivers two core components for AI CLI agents:
1. **OxHUD (`oxhud`)**: A lightning-fast, zero-dependency Go statusline HUD that monitors token context usage, active tools, and background agent tasks directly in your terminal.
2. **Skills & Plugins Ecosystem (`hub/`)**: A community registry of **120+ specialized skills** (Security, Coding, Refactoring, Research, MLOps, DevOps) with universal adapters for Claude Code, Codex, and Hermes.

---

## 🚀 Part 1: OxHUD (CLI Statusline)

OxHUD renders an elegant, single-line or two-line real-time status bar beneath or within your agent sessions.

```text
[CLAUDE: Opus 3.7] ██░░░░░░░░ 50k/200k (25.0%)
◐ Tool: edit_file | Task: patching state.go
```

### Installation

#### 1. Pre-built Binary or Build from Source

```bash
# Clone the repository
git clone https://github.com/TheOwlOps/OxHub.git
cd OxHub

# Build single binary
go build -o bin/oxhud ./cmd/oxhud

# (Optional) Install globally into your PATH:
# Linux/macOS:
sudo cp bin/oxhud /usr/local/bin/
# Windows (PowerShell as Admin):
Copy-Item bin/oxhud.exe C:\Windows\System32\
```

---

### Agent Integrations

#### 🤖 Claude Code
Configure Claude Code to invoke `oxhud` via its native statusline API:

1. Add configuration to your `~/.claude/settings.json`:
   ```json
   {
     "statusLine": "oxhud --agent=claude"
   }
   ```
2. Or test via stdin:
   ```bash
   echo '{"model":{"display_name":"Opus 3.7"},"context_window":{"current_usage":{"input_tokens":45000,"output_tokens":5000},"context_window_size":200000}}' | oxhud --agent=claude
   ```

#### 🦅 Hermes Agent
Install the `oxhud` skill or integrate directly into agent loops:

```bash
echo '{"model":"anti","used_tokens":85000,"total_tokens":256000,"active_tool":"terminal","task":"executing build"}' | oxhud --agent=hermes
```

#### 💻 OpenAI Codex
Pipe status events directly into `oxhud`:

```bash
echo '{"model":"gpt-4o","used_tokens":110000,"total_tokens":128000,"active_tool":"edit_file","status":"patching"}' | oxhud --agent=codex
```

---

## 📦 Part 2: Skills & Plugins Registry

OxHub includes over 120 production-ready skills indexed in `hub/manifest.json`.

### Featured Packs

| Category | Skills | Description |
| :--- | :--- | :--- |
| 🛡️ **Security** | `ox-security-audit`, `ox-pen-test`, `security-and-hardening` | Secret scanning, OWASP Top 10, Auth bypass, rate-limiting & network tests. |
| 💻 **Coding** | `ox-code-reviewer`, `ox-systematic-debug`, `ox-refactor-cleanup` | Concurrency leak checks, YAGNI enforcement, 4-phase root-cause debugging. |
| 🤖 **Agents** | `claude-code`, `hermes-agent`, `codex`, `computer-use` | Agent orchestration, multi-agent delegation, merge conflict reconciler. |
| 🧠 **MLOps** | `evaluating-llms-harness`, `llama-cpp`, `serving-llms-vllm`, `weights-and-biases` | GGUF quant, vLLM deployment, benchmark harnesses, W&B sweeps. |
| 📊 **Productivity**| `google-workspace`, `pdf`, `docx`, `xlsx`, `airtable`, `notion` | Full document pipelines, spreadsheet automation, headless office generation. |

### Using Skills

#### For Hermes Agent:
Copy any skill directory from `hub/skills/<category>/<skill-name>` to `~/.hermes/skills/`:
```bash
cp -r hub/skills/hub/security/security-audit ~/.hermes/skills/
```

#### For Claude Code & Codex:
Refer to the skill instruction file `SKILL.md` directly in prompts or include them in your `.claude/` or workspace guidelines:
```bash
claude --prompt "Follow guidelines in hub/skills/hub/coding/code-reviewer/SKILL.md"
```

---

## 🛠️ Development & Contributing

We welcome community contributions for new adapters, skills, and agents!

```bash
# Run tests
go test ./...

# Rebuild skills manifest
go run ./scripts/gen_manifest.go
```

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingSkill`)
3. Commit your Changes (`git commit -m 'Add AmazingSkill'`)
4. Push to the Branch (`git push origin feature/AmazingSkill`)
5. Open a Pull Request

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
