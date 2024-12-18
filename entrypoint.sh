#!/usr/bin/env bash

set -e
set -u
set -o pipefail

#if [ -n "${PARAMETER_STORE:-}" ]; then
#  export ESPACIOS_FISICOS_MID_PGUSER="$(aws ssm get-parameter --name /${PARAMETER_STORE}espacios_fisicos_mid/db/username --output text --query Parameter.Value)"
#  export ESPACIOS_FISICOS_MID_PGPASS="$(aws ssm get-parameter --with-decryption --name /${PARAMETER_STORE}/espacios_fisicos_mid/db/password --output text --query Parameter.Value)"

exec ./main "$@"
