---
name: ox-security-audit
description: Use when auditing source code for security vulnerabilities, secret leaks, and OWASP issues.
author: TheOwlOps
version: 1.0.0
---

# OX Security Audit & Hardening

Conduct systematic static and dynamic security assessments across Go, TypeScript, Python, and shell scripts.

## Core Checklist

1. **Secret & Credential Leaks**
   - High-entropy strings, private keys (`BEGIN RSA PRIVATE KEY`), JWT secrets, DB connection strings, AWS/GCP/OpenAI tokens.
   - Ensure `.gitignore` explicitly prevents `.env*`, `*.pem`, `*.key`, `credentials.json`.
   - Scan using `git log -S` or ripgrep regex patterns:
     ```bash
     rg -i "(api[_-]?key|secret|password|bearer|token)\s*[:=]\s*['\"][A-Za-z0-9_\-]{8,}['\"]"
     ```

2. **OWASP Top 10 & Input Sanitization**
   - **SQL Injection**: Enforce parameterized queries / ORM prepared statements. Zero string concatenation into SQL statements.
   - **Command Injection**: Avoid `exec.Command("bash", "-c", input)` or `os.system()` with unsanitized user arguments. Pass static command and argument array slice.
   - **SSRF**: Validate target IP/hostname against private subnet ranges (`10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`, `127.0.0.1/8`, `::1`).
   - **Path Traversal**: Always sanitize paths using `filepath.Clean` and verify target path has destination base prefix via `strings.HasPrefix(target, baseDir)`.

3. **Dependency Hardening**
   - Go: `go vet ./...` & `govulncheck ./...`
   - Node: `pnpm audit --audit-level=high` or `npm audit`
   - Python: `pip-audit` or `safety check`

4. **Output Format**
   - Severity: CRITICAL | HIGH | MEDIUM | LOW
   - File & Line Number: `path/to/file.ext:line`
   - Vulnerability Description & Exploit scenario
   - Concrete fix diff
