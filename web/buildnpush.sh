#!/usr/bin/env bash
set -euo pipefail
: "${IMAGE_REF:?Set IMAGE_REF to your own registry/image:tag}"
cd -- "$(dirname -- "$0")"
npm run build
docker build -t "$IMAGE_REF" .
docker push "$IMAGE_REF"
