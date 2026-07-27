# Security Policy — hawk-mcpkit

## Supported versions

We support the latest minor version on each `0.x` line, and the latest two
minor versions once `1.x` ships. Older versions receive critical-severity
fixes only on a best-effort basis.

The current canonical version is the contents of the [`VERSION`](./VERSION)
file at the repo root. See [`VERSIONING.md`](https://github.com/GrayCodeAI/hawk/blob/main/VERSIONING.md)
for the eco-wide versioning scheme.

## Reporting a vulnerability

**Do not open a public GitHub issue for security vulnerabilities.** Instead:

1. Open a private [GitHub Security Advisory](https://github.com/GrayCodeAI/hawk-mcpkit/security/advisories/new), **or**
2. Email `security@graycode.ai` with the details below.

Include in your report:

- A description of the vulnerability and the affected component.
- Steps to reproduce, ideally with a minimal proof-of-concept.
- The version (`VERSION` file or git SHA) you tested against.
- The potential impact and any suggested mitigation.

**Response targets:**

- Initial acknowledgement: within **48 hours**.
- Triage and severity assessment: within **5 business days**.
- Coordinated fix and disclosure: within **30 days** for high/critical, **90
  days** for medium/low (per industry-standard responsible disclosure).

## Disclosure policy

We follow [coordinated vulnerability disclosure](https://en.wikipedia.org/wiki/Coordinated_vulnerability_disclosure):

- Reporters receive credit in the advisory and CHANGELOG (unless they opt
  out).
- We request that reporters refrain from public disclosure until a fix has
  been released or the disclosure deadline above has elapsed.
- We will not pursue legal action against good-faith researchers acting
  within this policy.

## Security practices in this repo

- **Dependency monitoring:** vulnerable dependencies are detected by
  `govulncheck`, which runs on every CI build (see `.github/workflows/ci.yml`).
- **Static analysis:** `golangci-lint` is enforced in CI.
- **Vulnerability scanning:** `govulncheck` runs on every CI build.
- **Secret detection:** `trufflehog` runs on every CI build.
- **Reproducible builds:** the library is built with `go build`.
- **No secrets in source:** API keys are configuration, not constants.

## Scope

This policy covers the code in this repository and the release artefacts
published from it. It does not cover:

- Third-party dependencies (report to upstream).
- LLM provider services (report to the provider).
- Local filesystem misuse where an attacker already has shell access (out of
  threat model).

For hawk-mcpkit-specific threat-model notes, see the README and any docs in
this repo.
