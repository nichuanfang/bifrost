# Bifrost verification matrix

Run checks from the repository root unless a command says otherwise. Select the
smallest useful set from the files changed by the prepared upstream merge.

| Changed paths | Minimum check |
| --- | --- |
| `core/providers/<name>/` | `go test ./core/providers/<name>` when a root `go.work` exists, otherwise run it from `core` as `go test ./providers/<name>` |
| `core/`, `framework/`, or `transports/` Go code | Tests for the affected module; use `go test ./...` when the change is cross-cutting |
| `ui/` | `(cd ui && npm run typecheck)` |
| `.github/workflows/` | YAML/workflow lint if available; inspect action inputs and permissions |
| `transports/Dockerfile*`, Go build files, or UI build files | Docker Buildx smoke build when Docker is available |
| `plugins/pii-masking/` or `plugins/vision-extension/` | Run `go test ./...` from each module, inspect exported `.so` symbols, and build them with the same dynamic host/toolchain |

For this fork, provider response parsing changes should also keep the relevant
OpenAI-compatible, Gemini, and Mistral tests in the validation set. Do not
revert local provider fixes merely because upstream touched the same provider.

If a full test fails only because repository-provided MCP binaries, fixtures, or
external services are absent, report that as an environment limitation and run
the affected package tests separately. A compile failure, assertion failure, or
changed behavior remains a sync failure.
