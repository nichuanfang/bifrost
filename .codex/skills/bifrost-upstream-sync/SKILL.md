---
name: bifrost-upstream-sync
description: Safely merge upstream Bifrost changes into a fork while preserving fork commits, detecting conflicts, and validating affected code.
metadata:
  short-description: Safely sync Bifrost with upstream
---

# Bifrost Upstream Sync

Use this skill when maintaining the user's Bifrost fork and synchronizing it
with `https://github.com/maximhq/bifrost.git`. The default upstream branch is
discovered from `upstream/HEAD` and the fork remote is normally `origin`.

## Safety rules

- Work from the repository root and read `AGENTS.md`, this skill, and
  `references/fork-change-log.md` before changing code.
- Refuse to start when the worktree has staged, unstaged, or untracked files.
- Never use `git reset --hard`, `git checkout --`, force-push, or implicit rebase.
- Preserve fork commits with a non-rewriting merge.
- Create a timestamped backup branch before preparing the merge.
- If Git reports conflicts, list the files, abort the prepared merge, and stop for semantic review. Do not choose ours/theirs wholesale.
- A merge is not complete until affected subsystem checks pass.
- Do not push to `origin` unless the user explicitly requests it in the current request.

## Workflow

1. Resolve the repository root and current branch. Reject detached HEAD or a
   dirty worktree. Read `references/fork-change-log.md` so the current fork
   customizations are known before preparing a merge.
2. Run `scripts/sync_upstream.sh --prepare`. It verifies or adds the exact upstream remote, fetches with pruning, discovers upstream HEAD, creates a backup branch, and prepares `git merge --no-ff --no-commit`.
3. Inspect the prepared diff and compare it with
   `references/fork-change-log.md`. Confirm every recorded fork behavior remains
   present and identify any new fork-specific behavior that needs a ledger entry.
4. Read `references/verification.md` and run checks selected by the changed paths. Run targeted tests first and broader tests when practical.
5. If checks pass, create the merge commit with `git commit --no-edit`. If checks fail, run `git merge --abort`, retain the backup branch, and report the failing command and paths.
6. After a fork-specific behavior is added or changed, update
   `references/fork-change-log.md` with the date, behavior, affected paths,
   compatibility rule, and validation command. Keep this ledger focused on
   durable behavior rather than transient implementation details.
7. Push only when explicitly requested, using normal push semantics.

The helper is intentionally a prepare step so validation happens before the
merge commit exists. It is safe to rerun after an abort; each run creates a new
backup branch.
