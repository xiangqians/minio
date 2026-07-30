# pgsty/minio Security Advisories

This document summarizes fork-specific security fixes and closely related upgrade-impacting security notes in `pgsty/minio`. It is intentionally narrower than a full changelog and focuses on release-impacting security behavior.

## Advisories since `RELEASE.2026-03-21T00-00-00Z`

| ID | Fixed by | Affected area | Remote exploitability | Summary | Upgrade / workaround notes |
| :-- | :-- | :-- | :-- | :-- | :-- |
| `CVE-2026-33322` | `d24f449e0` | OIDC STS (`AssumeRoleWithWebIdentity`, `AssumeRoleWithClientGrants`) | Yes | Closes JWT algorithm confusion by removing HMAC/shared-secret verification and requiring JWKS-backed verifier keys | Breaking change: providers issuing `HS256`, `HS384`, or `HS512` tokens for these STS flows must switch to JWKS-backed RSA or ECDSA signing before upgrading. `PS256` and `EdDSA` are not currently supported. |
| `CVE-2026-33419` | `3b950f8fa` | LDAP STS authentication | Yes | Prevents username enumeration by unifying unknown-user and bad-password responses and adds in-memory login throttling | Unknown users and invalid passwords now both return `400 InvalidParameterValue`. Rate limiting is per-node, in-memory, and not currently configurable. Follow-up hardening landed in `18b712d49`, `9e10f6d9a`, and `f44110890`. |
| `CVE-2026-34204` | `56fa63bfd` | Replication metadata handling | Yes | Blocks untrusted `X-Minio-Replication-*` headers from being smuggled into internal replication metadata and leaving objects unreadable | Upgrade any server that accepts untrusted `PutObject` or `CopyObject` requests, which in practice means almost any production server that accepts writes. |
| `CVE-2026-39414` | `3252d5b7f` | S3 Select oversized record handling | Yes | Rejects oversized CSV records and non-`simdjson` line-delimited JSON records with `OverMaxRecordSize` instead of buffering them unchecked | Split oversized JSON lines client-side if you rely on `simdjson` input paths until the fast path gains the same pre-check. |
| `fake CVE-2026-40027` | `f444b6f37` | Unsigned-trailer PUT and multipart upload authentication | Yes | Closes the query-string authentication bypass in unsigned-trailer streaming requests | Upgrade if clients can reach object write endpoints using the `STREAMING-UNSIGNED-PAYLOAD-TRAILER` content-sha256 mode together with query-string SigV4 credentials. |
| `fake CVE-2026-40028` | `efb6e5b00` | Snowball auto-extract authentication | Yes | Verifies request authentication before tar extraction in Snowball unsigned-trailer flows | Upgrade if you use `PutObjectExtract` or Snowball uploads. |

## Dependency security updates

| ID | Fixed by | Summary |
| :-- | :-- | :-- |
| `CVE-2026-34986` | `68e0ba997` | Upgrades `go-jose` to `v4.1.4`. |
| `CVE-2026-39883` | `1869bd30b`, `e4fa06394` | Updates OpenTelemetry dependencies. |
| Upstream Go security fixes | `db4c0fd5e`, Go 1.26.4 update | Bumps the Go toolchain through `1.26.4`. |

## Operationally significant security-related fixes

| Change | Fixed by | Summary |
| :-- | :-- | :-- |
| LDAP TLS regression | `ce1c537eb` | Restores TLS configuration propagation for `ldaps://` `DialURL()` connections so `MINIO_IDENTITY_LDAP_TLS_SKIP_VERIFY` and custom root CAs work again. |
