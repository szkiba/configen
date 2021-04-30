#!/bin/sh
set -euo pipefail

if [ "${1#-}" != "${1}" ] || [ -z "$(command -v "${1}")" ]; then
  set -- configen "$@"
fi

exec "$@"
