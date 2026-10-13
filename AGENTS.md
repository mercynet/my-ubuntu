# Repository Guidelines

## Project Structure & Module Organization

Meu Ubuntu is a Go/Wails desktop application with React/TypeScript, targeting Ubuntu 26.04, GNOME, and amd64. `main.go` handles startup and embedded resources; `app.go` exposes the frontend bridge. Backend packages live in `internal/`, including planning, execution, detection, storage, and the privileged helper. Recipes are in `internal/catalog/catalog.json`; bundled profiles are in `internal/profiles/`.

The interface lives in `frontend/src/`, with local logos in `frontend/src/assets/apps/`. Go tests sit beside source files; frontend tests live in `frontend/tests/`. `build/` contains container, packaging, and verification scripts; `dist/` holds build outputs. Read `docs/desenvolvimento.md` for architectural contracts.

## Build, Test, and Development Commands

Use Docker and Make for the supported build environment:

- `make check` (or `make test`): check generated docs, install frontend dependencies, run Node tests, compile TypeScript/Vite, and run Go vet/tests with race detection.
- `make audit`: run Go and npm vulnerability audits.
- `make build`: produce binaries, a Debian package, and checksums in `dist/`.
- `make verify`: validate built checksums, versions, and package metadata.
- `make sim`: open the built application with the simulation executor.

For frontend changes, run `npm ci`, `npm test`, and `npm run build` inside `frontend/`. Build frontend assets before testing the complete Go application because Wails embeds `frontend/dist`.

## Coding Style & Naming Conventions

Format Go with `gofmt`; use lowercase package names and standard exported-name capitalization. Match existing frontend style: two-space indentation, double quotes, semicolons, PascalCase component files, and camelCase helpers. TypeScript uses strict checking; no dedicated frontend formatter or linter is configured. Keep recipe/profile IDs stable and lowercase with hyphens.

## Testing Guidelines

Use Go's `testing` package with `*_test.go` files and `Test...` functions; frontend tests use `node:test` in `*.test.mjs`. Add regression coverage for changed behavior, especially validation, conflicts, permissions, and persistence. No numeric coverage threshold is configured. Run `realsystem` tests only in disposable containers/VMs with `MEU_UBUNTU_DISPOSABLE=1`. Validate desktop behavior on GNOME when relevant.

## Commit & Pull Request Guidelines

Use concise, imperative commit subjects. Development happens on `development`; release tags use `v<version>`. Follow `RELEASING.md` for versions, evidence and publication. Generated files in `dist/` and `frontend/dist/` stay out of Git. PRs should explain behavior changes, link relevant issues, report validation and limitations, and include screenshots for interface changes.

## Security & Documentation

Preserve reviewed-plan tokens and helper operation allowlists; frontend requests carry IDs/options, never arbitrary shell commands. Exclude secrets from profiles and logs. Edit `docs/manual.json` and regenerate documentation with `python3 build/scripts/generate-docs.py`.
