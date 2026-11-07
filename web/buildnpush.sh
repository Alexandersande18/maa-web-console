#!/usr/bin/env bash
set -euo pipefail
cd -- "$(dirname -- "$0")"
if [[ -z "${IMAGE_REF:-}" && -f deploy.local.json ]]; then
    IMAGE_REF="$(node -e 'process.stdout.write(JSON.parse(require("fs").readFileSync("deploy.local.json", "utf8")).imageRef || "")')"
fi
: "${IMAGE_REF:?Set IMAGE_REF or configure the ignored deploy.local.json file}"
npm run build
docker build -t "$IMAGE_REF" .
docker push "$IMAGE_REF"
