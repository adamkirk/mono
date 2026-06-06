#!/bin/bash
set -e

BASE_SHA=$1
HEAD_SHA=$2

declare -A SCOPE_PATHS=(
  ["api"]="api/"
  ["ui"]="ui/"
  ["cli"]="cli/"
  ["chart-api"]="charts/api/"
  ["chart-ui"]="charts/ui/"
)

COMPONENT_DIRS=("api/" "ui/" "cli/" "charts/api/" "charts/ui/")

FAILED=0

while IFS= read -r commit; do
  message=$(git log --format="%s" -n 1 "$commit")

  regex='^\w+\(([^)]+)\):'
  if [[ "$message" =~ $regex ]]; then
    scope="${BASH_REMATCH[1]}"
  else
    echo "ERROR: Commit $commit has no scope: $message"
    FAILED=1
    continue
  fi

  files=$(git diff-tree --no-commit-id -r --name-only "$commit")

  if [ "$scope" = "repo" ]; then
    while IFS= read -r file; do
      for dir in "${COMPONENT_DIRS[@]}"; do
        if [[ "$file" == "$dir"* ]] && [[ "$(basename "$file")" != "release.config.cjs" ]] && [[ "$file" != *"/.githook/"* ]]; then
          echo "ERROR: Commit $commit (scope: repo) touches component file: $file"
          FAILED=1
        fi
      done
    done <<< "$files"
  else
    expected_path="${SCOPE_PATHS[$scope]}"

    if [ -z "$expected_path" ]; then
      echo "ERROR: Unknown scope '$scope' in commit: $message"
      FAILED=1
      continue
    fi

    while IFS= read -r file; do
      if [[ "$file" != "$expected_path"* ]]; then
        echo "ERROR: Commit $commit (scope: $scope) touches out-of-scope file: $file"
        FAILED=1
      fi
    done <<< "$files"
  fi
done < <(git log --format="%H" "$BASE_SHA..$HEAD_SHA")

if [ $FAILED -ne 0 ]; then
  exit 1
fi

echo "All commits pass scope validation"
