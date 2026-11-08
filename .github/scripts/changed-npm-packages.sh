#!/usr/bin/env bash
# changed-npm-packages.sh <commit>
#
# Prints the npm packages whose published contents can change between the
# newest release tag and the supplied commit. Package names are space-separated
# so the result can be passed through a GitHub Actions output.
set -euo pipefail

sha=$(git rev-parse "$1^{commit}")
newest=$(git tag --list 'v[0-9]*.[0-9]*.[0-9]*' --sort=-v:refname | awk 'NR == 1')

if [ -n "$newest" ] && ! git merge-base --is-ancestor "$newest" "$sha"; then
  exit 0
fi

cli=0
mcp=0

if [ -n "$newest" ]; then
  paths=$(git diff --name-only "$newest..$sha")
else
  paths=$(git log --format= --name-only "$sha")
fi

while IFS= read -r path; do
  case "$path" in
    packages/breeze-cli/*) cli=1 ;;
    packages/breeze-mcp/*|ui/breeze/catalog/*) mcp=1 ;;
  esac
done <<< "$paths"

packages=()
[ "$cli" -eq 1 ] && packages+=(breeze-cli)
[ "$mcp" -eq 1 ] && packages+=(breeze-mcp)
printf '%s\n' "${packages[*]-}"
