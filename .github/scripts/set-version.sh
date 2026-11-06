#!/usr/bin/env bash
# set-version.sh <vX.Y.Z>
#
# Writes a release's version into the files that carry one: the npm packages' package.json and the version the MCP server reports. The Go module and the breeze-ui command take theirs from the tag itself. Each file must hold its version exactly once, so that a moved or renamed line fails here rather than going unbumped.
set -euo pipefail

tag=$1
version=${tag#v}
if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "not a version: $tag" >&2
  exit 1
fi

root=$(git rev-parse --show-toplevel)

# bump <file> <pattern of the line> <the line with the new version>
bump() {
  local file=$root/$1 pattern=$2 line=$3 count
  count=$(grep -cE "$pattern" "$file" || true)
  if [ "$count" -ne 1 ]; then
    echo "$1: expected one line matching $pattern, found $count" >&2
    exit 1
  fi
  sed -E "s/$pattern/$line/" "$file" >"$file.tmp"
  mv "$file.tmp" "$file"
}

bump packages/breeze-cli/package.json '^  "version": "[^"]*",$' "  \"version\": \"$version\","
bump packages/breeze-mcp/package.json '^  "version": "[^"]*",$' "  \"version\": \"$version\","
bump packages/breeze-mcp/src/index.ts '^  version: "[^"]*",$' "  version: \"$version\","
