
## Commit conventions

All commits must follow [Conventional Commits](https://www.conventionalcommits.org/) format and include a **scope**:

```
<type>(<scope>): <description>

<body>
```

These rules are enforce in CI.

### Valid scopes

The scopes are required in the commit message to support the mono-repo structure. This repository publishes various packages, and each one requires separate versioning. Each of the scopes (aside from `repo`) are a component within this mono-repo.

| Scope | Component path |
|-------|---------------|
| `api` | `api/` |
| `ui` | `ui/` |
| `cli` | `cli/` |
| `chart-api` | `charts/api/` |
| `chart-ui` | `charts/ui/` |
| `repo` | Anything else. |

### Scope enforcement rules

- Each commit's scope must match the files it touches. A commit scoped to `api` may only modify files under `api/`, and so on.
- The `repo` scope is for repository-level changes, typically for development or release related changes. Commits with scope `repo` must **not** touch any component directory (`api/`, `ui/`, `cli/`, `charts/api/`, `charts/ui/`). The only exception is `release.config.cjs` files inside those directories, which are considered repo config.
- Commits that touch files in multiple components must be split into separate commits, one per scope.

## Releases

Releases are triggered automatically when a PR is merged to `main`. Each component is released independently based on which files changed.

### How versions are determined

Semantic versioning is driven by commit types for the relevant scope:

| Commit type | Version bump |
|-------------|-------------|
| `feat` | minor |
| `fix` | patch |
| `perf` | patch |
| `revert` | patch |

Commits scoped to a different component are ignored by that component's release pipeline, so merging changes to `api` will never trigger a `ui` release.

### Release tags

Each component is tagged independently:

| Component | Tag format |
|-----------|-----------|
| `api` | `api-v<version>` |
| `ui` | `ui-v<version>` |
| `cli` | `cli-v<version>` |
| `chart-api` | `chart-api-v<version>` |
| `chart-ui` | `chart-ui-v<version>` |
