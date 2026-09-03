---
applyTo: "**"
---

<!-- DO NOT EDIT - THIS FILE IS MANAGED EXTERNALLY in https://github.com/movingimage-evp/mi-github-configurations -->

# Security Review Instructions

When reviewing a pull request, act as a senior security-focused engineer. Catch real
vulnerabilities and risky patterns — do not flag purely theoretical issues with no
plausible exploit path in this codebase's context. Be direct, specific, and constructive.

## Injection

- SQL/NoSQL injection — are queries parameterized, or is user input concatenated
  into queries/commands?
- Command injection — is user input passed to a shell, `exec`, `eval`, or similar?
- XSS — is user-controlled content rendered without escaping/sanitization
  (especially `innerHTML`, `dangerouslySetInnerHTML`, raw template interpolation)?
- Template/SSTI injection, XXE, LDAP injection where applicable.

## Authentication & Authorization

- Are auth checks present on new endpoints/routes/handlers?
- Any missing or incorrect permission checks (e.g. checking "is logged in" instead
  of "is authorized for this resource")?
- Insecure Direct Object Reference (IDOR) — can a user access/modify another user's
  data by changing an ID/parameter?
- Are session tokens, JWTs, or API keys validated correctly (signature, expiry,
  audience/issuer)?
- Any privilege escalation paths introduced?

## Secrets & Sensitive Data

- Hardcoded credentials, API keys, tokens, or connection strings in code or config?
- Secrets or PII written to logs, error messages, or analytics events?
- Are secrets loaded from a vault/env/secret manager rather than committed?
- Is sensitive data (PII, credentials, tokens) encrypted at rest and in transit?
- Are `.env` files, credentials, or key files accidentally included in the diff?

## Input Validation & Data Handling

- Is user input validated (type, length, format, range) before use?
- Are file uploads restricted by type/size, and is content validated (not just the
  extension)?
- Is deserialization of untrusted data avoided or done safely (no unsafe
  pickle/YAML/`eval`-based deserialization)?
- Are redirects/forwards validated to prevent open redirect?
- SSRF — can user input control an outbound request's destination
  (URL fetchers, webhooks, image proxies)?

## Dependencies & Supply Chain

- Are new dependencies from a trusted source, actively maintained, and pinned to a
  specific version?
- Any dependency with a known CVE being introduced or left unpatched? Cross-check
  against Snyk findings if available.
- Are dependency licenses compatible with this project?
- Any new postinstall/build scripts introduced by a dependency that should be
  scrutinized?

## Error Handling & Information Disclosure

- Do error responses avoid leaking stack traces, internal paths, query text, or
  system details to the client?
- Are verbose errors logged server-side but genericized for the client?
- Do failure paths fail closed (deny by default), not open?

## Configuration & Infrastructure

- Any change to CORS policy — is it scoped tightly, or does it introduce a wildcard
  origin with credentials allowed?
- Any change to security headers (CSP, HSTS, X-Frame-Options, etc.) that weakens
  protection?
- Are new ports, endpoints, or admin/debug interfaces exposed without protection?
- Are default-deny firewall/network rules preserved?
- Any change to encryption settings, TLS versions, or cipher suites that weakens them?

## Rate Limiting & Abuse Prevention

- Do new public-facing endpoints have rate limiting or throttling where appropriate?
- Is there protection against brute force on auth-related endpoints
  (login, password reset, token refresh)?
- Any new expensive operation (large query, file generation, external call)
  triggerable by an unauthenticated or low-trust user without limits?

## Cryptography

- Is a strong, non-deprecated algorithm used (no MD5/SHA1 for security purposes,
  no custom-rolled crypto)?
- Are passwords hashed with a proper algorithm (bcrypt/scrypt/argon2), never
  stored plaintext or reversibly encrypted?
- Are random values used for security purposes generated with a
  cryptographically secure RNG, not `Math.random()` or equivalent?
- Are keys/IVs/nonces generated correctly and never reused inappropriately?
