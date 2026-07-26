#!/usr/bin/env bash
# Regenerates PROGRESS.md from the directory tree, keeping any boxes you ticked.
set -euo pipefail
cd "$(dirname "$0")"
out=PROGRESS.md
old=""
[[ -f $out ]] && old=$(cat "$out")

{
  echo "# Progress"
  echo
  echo "Tick a box when \`go test\` is green for that exercise."
  echo "Regenerate with \`./gen_progress.sh\` — your ticks are preserved."
  echo
  total=0
  for mod in [0-9][0-9]-*/; do
    echo "## ${mod%/}"
    echo
    for ex in "$mod"*/; do
      [[ -d $ex ]] || continue
      name="${ex%/}"
      box=" "
      if grep -qF -- "[x] $name" <<<"$old"; then box="x"; fi
      echo "- [$box] $name"
      total=$((total+1))
    done
    echo
  done
  echo "---"
  echo
  echo "$total items."
} > "$out"
echo "wrote $out"
