# Bifrost Fork Change Log

This ledger records durable behavior that exists in the `nichuanfang/bifrost`
fork and must be checked after upstream synchronization. It replaces the need
to encode fork-specific behavior in the repository's general `AGENTS.md`.

## 2026-09-13

### Provider model-list response decoding

- Commit: `fix(providers): decode compressed model list responses`
- Behavior: model-list responses are passed through
  `providerUtils.CheckAndDecodeBody` before JSON parsing. Ordinary uncompressed
  JSON remains unchanged; compressed responses receive correct decoding and
  malformed encodings produce a distinct provider response decode error.
- Affected paths: `core/providers/anthropic`, `elevenlabs`, `gemini`, `mistral`,
  `openai`, `openrouter`, `replicate`, and `vertex`.
- Compatibility rule: preserve the decode step when resolving upstream changes
  in shared OpenAI-compatible or provider-specific model-list code.
- Validation: targeted tests for the affected provider packages, including
  OpenAI, Gemini, and Mistral compression coverage.

### OSS enterprise feature boundary

- Commit: `fix(ui): hide enterprise features and remove release prompt`
- Behavior: OSS builds hide enterprise-only routes and UI sections, reject
  direct navigation to enterprise-only routes, and keep enterprise-only prompt
  and observability components out of the OSS UI.
- Affected paths: `ui/app/workspace/layout.tsx`,
  `ui/components/sidebar.tsx`,
  `ui/lib/constants/enterprise.ts`, and related feature views.
- Compatibility rule: classify newly added enterprise routes and extend the
  route boundary before allowing them into the OSS navigation.
- Validation: UI typecheck and enterprise-route unit tests.

### Web UI release prompt removal

- Commit: `fix(ui): hide enterprise features and remove release prompt`
- Behavior: the Web UI no longer queries `https://getbifrost.ai/latest-release`
  or renders the release notification card. CLI self-update behavior is
  unchanged.
- Affected paths: `ui/components/sidebar.tsx`,
  `ui/lib/store/apis/configApi.ts`, and `ui/lib/types/config.ts`.
- Compatibility rule: do not reintroduce the Web UI release query or card
  unless the fork explicitly opts back into update notifications.
- Validation: UI typecheck and targeted lint.

### GHCR amd64 publishing from local modules

- Commit: `ci: publish amd64 GHCR images from local modules`
- Behavior: tag and manual GHCR publishes build with
  `transports/Dockerfile.local` and publish only `linux/amd64`, so the image
  contains the checked-out fork versions of core, framework, plugins, and
  transports.
- Affected paths: `.github/workflows/publish-ghcr.yml` and
  `transports/go.sum`.
- Compatibility rule: keep the local-module Dockerfile, GHCR write
  permission, metadata tags, `VERSION` build argument, and amd64-only target.
  Do not silently switch to the proxy-based `transports/Dockerfile`.
- Validation: workflow inspection, transports dependency validation, and a
  Docker Buildx smoke build when Docker is available.

### Custom PII and vision plugins in the local-module image

- Commit: `feat(plugins): integrate custom pii and vision plugins`
- Behavior: the fork includes the `pii-masking` and `vision-extension` native
  plugins. `transports/Dockerfile.local` builds the Bifrost host dynamically
  and compiles both `.so` files from the same local Go workspace, toolchain,
  libc implementation, and target architecture before copying them into
  `/app/plugins/`. Plugins remain opt-in through explicit `config.json`
  entries; the standard `transports/Dockerfile` remains the static image.
- Affected paths: `plugins/pii-masking`, `plugins/vision-extension`,
  `transports/Dockerfile.local`, and
  `examples/configs/withcustomplugins/config.json`.
- Compatibility rule: whenever core, Go, the Alpine/musl toolchain, or the
  target architecture changes, rebuild both plugins together with the host
  binary. Preserve dynamic linking for the plugin-enabled image and keep
  plugin paths under `/app/plugins/`.
- Validation: run `go test ./...` in both custom plugin modules, run the
  dynamic-plugin loader tests in `framework/plugins`, inspect exported plugin
  symbols, validate the example JSON, and run a Docker Buildx smoke build
  before publishing the image.
