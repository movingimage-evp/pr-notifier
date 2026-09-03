---
applyTo: "**/Dockerfile*"
---

<!-- DO NOT EDIT - THIS FILE IS MANAGED EXTERNALLY in https://github.com/movingimage-evp/mi-github-configurations -->

# Docker Review Instructions

## Purpose

When reviewing changes to a `Dockerfile` file, or related container
config, act as a senior engineer focused on image security, size, and build
reliability. Catch real issues — do not flag purely stylistic preferences with no
practical impact. Be direct, specific, and constructive.

## Base Image

- Is the base image pinned to a specific version tag, not `latest`?
- Better: is it pinned by digest (`@sha256:...`) for full reproducibility, where the
  project's maturity warrants it?
- Is the base image from a trusted/official source (Docker Official Images,
  verified publisher), not an unmaintained or unofficial one?
- Is a minimal base used where possible (`-slim`, `-alpine`, or `distroless`)
  instead of a full OS image, unless there's a documented reason not to
  (e.g. glibc dependency, native module compatibility)?
- Is the base image reasonably up to date (not years-stale with known CVEs)?

## Multi-Stage Builds

- For compiled/bundled languages, is a multi-stage build used so build tools,
  source, and dev dependencies don't end up in the final image?
- Does the final stage copy only the build artifacts it needs
  (`COPY --from=build ...`), not the entire build context?
- Are build-time secrets (private registry tokens, SSH keys) kept out of the final
  image layers, ideally via `--mount=type=secret` (BuildKit) rather than `ARG`/`ENV`?

## Layer Caching & Build Efficiency

- Are dependency manifests (`package.json`/`requirements.txt`/`go.mod`/etc.)
  copied and installed in a separate, earlier layer from the application source,
  so dependency install is cached across source-only changes?
- Are related `RUN` commands chained with `&&` to avoid unnecessary intermediate
  layers, and are package manager caches cleaned in the same layer they're created
  (e.g. `rm -rf /var/lib/apt/lists/*` in the same `RUN` as `apt-get install`)?
- Is there a `.dockerignore` file, and does it exclude `.git`, `node_modules`,
  local env files, build artifacts, and other content not needed in the image?

## User & Permissions

- Does the container run as a non-root user in the final image
  (`USER` instruction set, not left as default root)?
- If a non-root user is created, is it done with a fixed, explicit UID/GID
  (helps with volume permission consistency), rather than an auto-assigned one?
- Are file permissions on copied application files no more permissive than needed?
- If the app needs to bind a privileged port (<1024), is that handled via a
  reverse proxy / port remap instead of running the container as root?

## Secrets & Sensitive Data

- No secrets, credentials, API keys, or private keys hardcoded in `ENV`, `ARG`,
  or `COPY`'d files — `ARG` values are visible in image history even if not in
  the final layer.
- No `.env` files, credentials, or key material copied into the image by an
  overly broad `COPY . .` — check that `.dockerignore` covers this.
- Are secrets injected at runtime (orchestrator secrets, mounted files, env vars
  from a secret manager) rather than baked into the image?

## Networking & Exposed Surface

- Does `EXPOSE` (and any `docker-compose` port mapping) match only the ports the
  service actually needs to expose?
- Is the container's attack surface minimized — no unnecessary packages, debug
  tools, compilers, or shells left in the final production image?
- In `docker-compose`, are internal-only services (databases, caches) kept off
  host-published ports, or bound to `127.0.0.1` rather than `0.0.0.0`, in
  non-dev configurations?

## Reliability & Operations

- Is a `HEALTHCHECK` defined (or configured at the orchestration layer) so the
  container's health is observable?
- Does the app handle `SIGTERM` for graceful shutdown, and does the Dockerfile
  avoid patterns that break signal propagation (e.g. shell-form `CMD` wrapping
  the process, missing an init system like `tini` for zombie reaping when
  needed)?
- Is `CMD`/`ENTRYPOINT` using exec form (`["executable", "arg"]`) rather than
  shell form, so signals reach the process directly and it runs as PID 1
  correctly?
- Are resource limits (CPU/memory) set at the orchestration/compose level where
  this environment expects them?

## Reproducibility

- Are dependency versions pinned (lockfiles committed and used, package versions
  specified) so builds are reproducible rather than pulling latest at build time?
- Does the build avoid non-deterministic steps (e.g. pulling `HEAD` of a branch,
  unpinned `curl | sh` installs) without integrity verification (checksum/signature)?

