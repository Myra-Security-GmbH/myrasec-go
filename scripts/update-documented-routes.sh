#!/usr/bin/env bash
set -euo pipefail

for tool in curl jq; do
    command -v "$tool" >/dev/null || { echo "$tool is required" >&2; exit 1; }
done

host="${1:-https://apiv2.myracloud.com}"
while [[ "$host" == */ ]]; do
    host="${host%/}"
done

output="$(cd "$(dirname "$0")/.." && pwd)/testdata/documented-routes.txt"
tmp="$(mktemp "${output}.XXXXXX")"
trap 'rm -f "$tmp"' EXIT

curl -fsSL --proto '=https' --proto-redir '=https' --max-time 30 --max-redirs 3 "${host}/api/doc.json" \
    | jq -r '
        if (.basePath // "/") != "/" then error("unexpected basePath \(.basePath)") else . end
        | .paths | to_entries[] | .key as $path
        | if ($path | test("^/[A-Za-z0-9_./{}-]*$")) then . else error("unexpected path \($path)") end
        | .value | keys[] | select(test("^(get|post|put|patch|delete)$")) | "\(ascii_upcase) \($path)"' \
    | sed -E 's/\{[^}]+\}/*/g' \
    | LC_ALL=C sort -u > "$tmp"

if [[ ! -s "$tmp" ]]; then
    echo "${host}/api/doc.json lists no operations, keeping testdata/documented-routes.txt" >&2
    exit 1
fi

mv "$tmp" "$output"
trap - EXIT

echo "Wrote $(wc -l < "$output") documented operations from ${host}/api/doc.json to testdata/documented-routes.txt"
