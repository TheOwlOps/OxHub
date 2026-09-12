# OxHub ⚡

> **Universal Agent Hub & Ecosystem** — High-performance statusline HUD, skills registry, and plugin ecosystem for **Claude Code**, **Hermes Agent**, and **OpenAI Codex**.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://golang.org)
[![Cross-Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-brightgreen)](#)

---

## ⚡ 1-Line Universal Install (Community & Agents)

Cài toàn bộ 120+ Skills vào agent của bạn bằng 1 lệnh duy nhất:

### Linux / macOS:
```bash
curl -fsSL https://raw.githubusercontent.com/TheOwlOps/OxHub/main/install.sh | bash -s -- all
```
*(Hoặc chọn từng agent: `bash -s -- hermes`, `bash -s -- claude`, `bash -s -- codex`)*

### Windows (PowerShell):
```powershell
irm https://raw.githubusercontent.com/TheOwlOps/OxHub/main/install.ps1 | iex
```

---

## 🎯 Native CLI Commands (Cho từng Agent)

### 1. Hermes Agent
Hermes hỗ trợ quản lý qua hệ thống `skills tap` và `skills install`:

- **Add toàn bộ repository OxHub làm source**:
  ```bash
  hermes skills tap add TheOwlOps/OxHub
  ```
- **Install lẻ từng skill qua URL/ID**:
  ```bash
  hermes skills install TheOwlOps/OxHub/hub/security/security-audit
  # Hoặc URL trực tiếp:
  hermes skills install https://raw.githubusercontent.com/TheOwlOps/OxHub/main/hub/security/security-audit/SKILL.md
  ```

---

### 2. Claude Code
- **Cài qua Plugin Manager**:
  ```bash
  /plugin marketplace add TheOwlOps/OxHub
  /plugin install oxhud
  ```
- **Auto StatusLine setup**:
  Thêm vào file cấu hình `~/.claude/settings.json`:
  ```json
  {
    "statusLine": "oxhud --agent=claude"
  }
  ```
- **Khai thác Skills**:
  Toàn bộ skills tải về thư mục `~/.claude/skills/` sẽ được Claude Code tự động nhận diện trong mọi phiên làm việc.

---

### 3. OpenAI Codex
- Chạy wrapper trực tiếp hoặc load skill từ thư mục `~/.codex/skills/`:
  ```bash
  echo '{"model":"gpt-4o","used_tokens":110000,"total_tokens":128000,"active_tool":"edit_file","status":"patching"}' | oxhud --agent=codex
  ```

---

## 🚀 OxHUD (CLI Real-time Statusline)

Statusline viết bằng Go, zero-dependency, siêu nhẹ:

```text
[CLAUDE: Opus 3.7] ██░░░░░░░░ 50k/200k (25.0%)
◐ Tool: edit_file | Task: patching state.go
```

```bash
# Build binary
go build -o bin/oxhud ./cmd/oxhud
```

---

## 📦 Skills & Plugins Catalog (120+ Skills)

Tất cả kỹ năng được lập chỉ mục tại `hub/manifest.json`.

| Nhóm | Skills Nổi Bật | Công Dụng |
| :--- | :--- | :--- |
| 🛡️ **Security** | `ox-security-audit`, `ox-pen-test`, `security-and-hardening` | Rà quét secret, phòng chống SQLi/RCE/SSRF, kiểm tra lỗ hổng OWASP & API. |
| 💻 **Coding** | `ox-code-reviewer`, `ox-systematic-debug`, `ox-refactor-cleanup` | Review code đa trục, YAGNI, 4-phase debugging, dọn dẹp mã nguồn thừa. |
| 🤖 **Autonomous** | `claude-code`, `hermes-agent`, `codex`, `computer-use` | Điều phối multi-agent, quản lý worktree, tự động hoá desktop/browser. |
| 🧠 **MLOps** | `evaluating-llms-harness`, `llama-cpp`, `serving-llms-vllm`, `weights-and-biases` | Quantize GGUF, triển khai vLLM, đánh giá benchmark, W&B sweep. |
| 📊 **Productivity**| `google-workspace`, `pdf`, `docx`, `xlsx`, `airtable`, `notion` | Tự động hoá tài liệu, xử lý bảng tính, OCR và tương tác workspace. |

---

## 📄 License
Phát hành theo giấy phép [MIT](LICENSE).
