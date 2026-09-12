---
name: pen-test
description: Use when assessing API endpoints, authentication flows, and network services for penetration security risks.
author: TheOwlOps
version: 1.0.0
---

# API & Network Penetration Testing

Assess attack surfaces on REST APIs, WebSocket channels, and RPC servers.

## Verification Checklist

1. **Authentication & Authorization Bypass**
   - Test missing `Authorization` header against protected routes.
   - BOLA / IDOR: Check if User A can mutate or query resources of User B by modifying path IDs (`/api/users/{id}/order`).
   - Token revocation: Confirm invalid/expired tokens are rejected immediately (HTTP 401).

2. **Rate Limiting & Abuse Prevention**
   - Burst test endpoints without auth or public endpoints (`/login`, `/signup`, `/search`).
   - Ensure rate limit headers (`Retry-After`, `X-RateLimit-Remaining`) return properly upon limit breach (HTTP 429).

3. **CORS & Security Headers**
   - Check `Access-Control-Allow-Origin` is not wildcard `*` with credentials enabled (`Access-Control-Allow-Credentials: true`).
   - Validate security headers:
     - `X-Content-Type-Options: nosniff`
     - `X-Frame-Options: DENY`
     - `Content-Security-Policy`
     - `Strict-Transport-Security`

4. **Payload Injection & Fuzzing**
   - Test extreme JSON depth, huge string buffers (>10MB) to check memory exhaustion (DoS).
   - Test malformed UTF-8 byte sequences and null byte injections.
