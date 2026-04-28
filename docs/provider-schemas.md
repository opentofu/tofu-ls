# Provider Schemas

Completion, hover, and validation for `resource` and `data` blocks come from
a provider schema. The version of that schema determines which attributes
are considered valid, so a stale schema can produce false `Unexpected
attribute` diagnostics.

## Sources

The server uses one of two sources for each provider:

**Bundled.** A selection of popular providers is embedded into every
`tofu-ls` release, one version per provider, refreshed at release time. Used
when no local schema is available, typically before `tofu init` has been
run.

**Local (from `tofu init`).** When a directory contains `.terraform/` and
`.terraform.lock.hcl`, the server runs `tofu providers schema -json` and
uses the result. This requires `tofu` on `PATH`. Re-run when the lock file
or module manifest changes.

## How the server picks between them

When more than one schema exists for the same provider, candidates are
scored and the highest wins:

| Factor | Score |
| --- | --- |
| `tofu init` schema in the current module | +2 |
| `tofu init` schema in a different module in the workspace | 0 |
| Bundled schema | -1 |
| Version satisfies the module's `required_providers` constraint | +2 |

A local schema for the current module always beats the bundled one. The
bundled schema is the fallback.

## Troubleshooting

If `Unexpected attribute` is reported for an attribute added in a newer
provider release:

- If you have run `tofu init`, the lock file is pinning an older version.
  Run `tofu init -upgrade`; the server picks up the change when the lock
  file is updated.
- If you have not run `tofu init`, run it. If the bundled schema is also
  older than the feature, update `tofu-ls` (or the editor extension that
  ships it).

If the diagnostic still seems wrong, you can disable schema-based
diagnostics with
[`validation.enableEnhancedValidation`](./SETTINGS.md#enablxxeenhancedvalidation-bool-defaults-to-true)
and [open an issue](https://github.com/opentofu/tofu-ls/issues/new/choose).
