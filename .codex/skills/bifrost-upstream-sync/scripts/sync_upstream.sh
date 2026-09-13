#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: sync_upstream.sh [options]

Prepare a safe, uncommitted merge of upstream Bifrost changes.

Options:
  --repo PATH             Repository path (default: current directory)
  --remote NAME           Upstream remote name (default: upstream)
  --url URL               Expected upstream URL
  --branch NAME           Upstream branch (default: upstream/HEAD, then dev)
  --prepare               Prepare the merge without committing (default)
  --dry-run               Fetch/check only; do not add a remote or create a branch
  --push                  Push the current branch after a prepared merge has been committed
  -h, --help              Show this help
EOF
}

repo="."
remote="upstream"
expected_url="https://github.com/maximhq/bifrost.git"
upstream_branch=""
dry_run=0
push=0

while (($# > 0)); do
  case "$1" in
    --repo) repo="$2"; shift 2 ;;
    --remote) remote="$2"; shift 2 ;;
    --url) expected_url="$2"; shift 2 ;;
    --branch) upstream_branch="$2"; shift 2 ;;
    --prepare) shift ;;
    --dry-run) dry_run=1; shift ;;
    --push) push=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; usage >&2; exit 2 ;;
  esac
done

repo="$(cd "$repo" && git rev-parse --show-toplevel)"
cd "$repo"

current_branch="$(git branch --show-current)"
if [[ -z "$current_branch" ]]; then
  echo "Refusing to sync from detached HEAD." >&2
  exit 1
fi

if ! git diff --quiet || ! git diff --cached --quiet || [[ -n "$(git ls-files --others --exclude-standard)" ]]; then
  echo "Refusing to sync a dirty worktree. Commit or intentionally stash the changes first." >&2
  git status --short >&2
  exit 1
fi

if git remote get-url "$remote" >/dev/null 2>&1; then
  actual_url="$(git remote get-url "$remote")"
  if [[ "$actual_url" != "$expected_url" ]]; then
    echo "Remote '$remote' points to '$actual_url', expected '$expected_url'." >&2
    echo "Refusing to overwrite the remote automatically." >&2
    exit 1
  fi
elif ((dry_run)); then
  echo "Remote '$remote' is absent; dry-run would add '$expected_url'."
  exit 0
else
  git remote add "$remote" "$expected_url"
  echo "Added remote '$remote' -> $expected_url"
fi

if ((dry_run)); then
  git fetch "$remote" --prune
  if [[ -z "$upstream_branch" ]]; then
    upstream_branch="$(git symbolic-ref --quiet --short "refs/remotes/$remote/HEAD" 2>/dev/null || true)"
    upstream_branch="${upstream_branch#"$remote/"}"
    upstream_branch="${upstream_branch:-dev}"
  fi
  echo "Dry-run: upstream branch is $upstream_branch"
  echo "Commits on upstream not in the current branch:"
  git log --oneline "HEAD..$remote/$upstream_branch" 2>/dev/null || true
  exit 0
fi

git fetch "$remote" --prune

if [[ -z "$upstream_branch" ]]; then
  upstream_branch="$(git symbolic-ref --quiet --short "refs/remotes/$remote/HEAD" 2>/dev/null || true)"
  upstream_branch="${upstream_branch#"$remote/"}"
  upstream_branch="${upstream_branch:-dev}"
fi

upstream_ref="refs/remotes/$remote/$upstream_branch"
if ! git rev-parse --verify "$upstream_ref" >/dev/null 2>&1; then
  echo "Upstream branch '$upstream_branch' was not found after fetch." >&2
  exit 1
fi

backup_branch="backup/upstream-sync-${current_branch//\//-}-$(date -u +%Y%m%d-%H%M%S)"
git branch "$backup_branch" HEAD
echo "Created backup branch: $backup_branch"
echo "Preparing merge: $remote/$upstream_branch -> $current_branch"

if ! git merge --no-ff --no-commit "$upstream_ref"; then
  conflicts="$(git diff --name-only --diff-filter=U || true)"
  git merge --abort >/dev/null 2>&1 || true
  echo "Upstream merge has conflicts; the prepared merge was aborted." >&2
  if [[ -n "$conflicts" ]]; then
    echo "Conflict files:" >&2
    printf '%s\n' "$conflicts" >&2
  fi
  echo "Backup branch retained: $backup_branch" >&2
  exit 2
fi

if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "Merge prepared but not committed. Run validation before: git commit --no-edit"
  git status --short
else
  echo "Already up to date; no merge commit is required."
fi

if ((push)); then
  if git rev-parse -q --verify MERGE_HEAD >/dev/null 2>&1; then
    echo "Refusing to push while the merge is uncommitted." >&2
    exit 1
  fi
  git push origin "$current_branch"
fi
