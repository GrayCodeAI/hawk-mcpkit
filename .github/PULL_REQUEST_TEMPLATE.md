<!--
  Thanks for your contribution! Please fill out this template so reviewers can
  understand the change quickly. Anything that does not apply can be left in
  place; do not delete unanswered sections — write "n/a".
-->

## Summary

<!--
  One paragraph describing what this PR does and why. Link the related
  issue(s) with `Fixes #N` or `Refs #N` if applicable.
-->

## Changes

<!--
  Bullet list of what changed, grouped by area (server, transport, auth,
  helpers, types, version, CI, docs).
  Reviewers should be able to skim this and know what to look at first.
-->

-

## API impact

<!--
  Did you add, remove, rename, or change the signature of any exported
  symbol? List them here. If yes, confirm whether this is a breaking
  change and bump the version accordingly in `VERSION`. If no exported
  surface changed, write "n/a".
-->

## Ecosystem boundary check

<!--
  hawk-mcpkit must never import hawk-eco dependencies. Confirm this PR
  does not add any imports of `hawk`, `eyrie`, `yaad`, `tok`, `trace`,
  `sight`, `inspect`, `hawk-core-contracts`, or any SDK.
-->

## Testing

<!--
  Describe how you tested. Paste output of `make test` and `make lint`.
  If you added new tests, list them.
-->

```text
$ make test
...
$ make lint
...
$ make boundaries
...
```

## Checklist

- [ ] Commits follow [Conventional Commits](https://www.conventionalcommits.org/)
      (`feat:`, `fix:`, `perf:`, `refactor:`, `docs:`, `test:`, etc.)
- [ ] `make ci` passes locally
- [ ] `make boundaries` passes (no hawk-eco imports)
- [ ] New or changed code has tests
- [ ] Public APIs have doc comments
- [ ] `CHANGELOG.md` updated under `## [Unreleased]` if user-visible
- [ ] `VERSION` file is **not** edited manually — release-please bumps it
- [ ] No secrets, tokens, or PII added to the repo
