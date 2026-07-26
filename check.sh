#!/usr/bin/env bash
# check.sh — show how far you've got.
#   ./check.sh          run every package, print a pass/fail summary
#   ./check.sh 06       run only module 06
#   ./check.sh -v 06    same, but show the raw go test output
set -uo pipefail
cd "$(dirname "$0")"

verbose=0
if [[ "${1:-}" == "-v" ]]; then verbose=1; shift; fi
filter="${1:-}"

pattern="./..."
if [[ -n "$filter" ]]; then pattern="./${filter}*/..."; fi

if [[ $verbose == 1 ]]; then
  go test -race -timeout 90s "$pattern"
  exit $?
fi

tmp=$(mktemp)
go test -race -count=1 -timeout 90s "$pattern" >"$tmp" 2>&1
pass=0; fail=0; noimpl=0
while IFS= read -r line; do
  case "$line" in
    ok*)    pass=$((pass+1)) ;;
    FAIL*goex/*) fail=$((fail+1)) ;;
  esac
done < "$tmp"
noimpl=$(grep -c 'TODO' "$tmp")
printf '\n  \033[32mpassing packages: %3d\033[0m\n' "$pass"
printf '  \033[31mfailing packages: %3d\033[0m\n' "$fail"
printf '  (%d tests still hitting a TODO stub)\n\n' "$noimpl"
echo "  failing:"
grep '^FAIL' "$tmp" | sed 's|^FAIL|   -|' | sed 's|\[build failed\]|(does not compile)|' | head -60
rm -f "$tmp"
