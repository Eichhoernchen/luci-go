#!/bin/bash
# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

METADATA_URL="http://metadata.google.internal"

_help() {
  echo -e "Usage: $0 <logdog-coordinator-path>"
  echo -e ""
  echo -e "  logdog-coordinator-path\tThe path to the application binary."
}

APP=$1; shift
if [ -z "${APP}" ]; then
  echo "ERROR: Missing argument <logdog-coordinator-path>"
  _help
  exit 1
fi
APP=$(readlink -f "${APP}")
if [ ! -x "${APP}" ]; then
  echo "ERROR: Application path is not an executable file [${APP}]"
  exit 1
fi

_load_metadata_domain() {
  local __resultvar=$1; shift
  local domain=$1; shift
  local resource=$1; shift

  OUTPUT=$(curl \
    -s \
    --fail \
    "${METADATA_URL}/computeMetadata/v1/${domain}/attributes/${resource}" \
    -H "Metadata-Flavor: Google")
  RV=$?
  if [ ${RV} != 0 ]; then
    return ${RV}
  fi
  if [ -z "${OUTPUT}" ]; then
    return 1
  fi

  eval $__resultvar="'${OUTPUT}'"
}

# Loads metadata from the GCE metadata server.
_load_metadata() {
  local __resultvar=$1; shift
  local resource=$1; shift

  # Try loading from 'instance' domain.
  _load_metadata_domain "$__resultvar" "instance" "$resource"
  if [ $? == 0 ]; then
    return 0
  fi

  # Try loading from 'project' domain.
  _load_metadata_domain "$__resultvar" "project" "$resource"
  if [ $? == 0 ]; then
    return 0
  fi

  echo "WARN: Failed to load metadata [${resource}]."
  return 1
}

_load_metadata_check() {
  local __resultvar=$1; shift
  local resource=$1; shift

  _load_metadata "$__resultvar" "$resource"
  if [ $? != 0 ]; then
    echo "ERROR: Metadata resource [${resource}] is required."
    exit 1
  fi

  return 0
}

_write_credentials() {
  local path=$1; shift
  local data=$1; shift
  echo "${data}" > "${path}"
}

# Test if we're running on a GCE instance.
curl \
  -s \
  --fail \
  "${METADATA_URL}" \
  1>/dev/null
if [ $? != 0 ]; then
  echo "ERROR: Not running on GCE instance."
  exit 1
fi

# Get the number of CPUs.
#
# TODO(dnj): Remove this when we switch to Go 1.5+.
NCPU=$(cat /proc/cpuinfo | grep processor | wc -l)
export GOMAXPROCS="${NCPU}"

# Load metadata (common).
_load_metadata_check COORDINATOR_HOST "logdog_coordinator_host"
_load_metadata STORAGE_CREDENTIALS "logdog_storage_auth_json"

# Load metadata (app-specific).
_load_metadata LOG_LEVEL "logdog_archivist_log_level"

# Runtime temporary directory.
TEMPDIR=$(mktemp -d)
trap "rm -rf ${TEMPDIR}" EXIT

# Compose command line.
ARGS=(
  "-coordinator" "${COORDINATOR_HOST}"
  "-config-kill-interval" "5m"
  )

if [ ! -z "${LOG_LEVEL}" ]; then
  ARGS+=("-log-level" "${LOG_LEVEL}")
fi
if [ ! -z "${STORAGE_CREDENTIALS}" ]; then
  STORAGE_CREDENTIALS_JSON_PATH="${TEMPDIR}/storage_service_account_json.json"
  _write_credentials "${STORAGE_CREDENTIALS_JSON_PATH}" "${STORAGE_CREDENTIALS}"
  ARGS+=("-storage-credential-json-path" "${STORAGE_CREDENTIALS_JSON_PATH}")
fi

echo "INFO: Running command line args: ${APP} ${ARGS[*]}"
"${APP}" ${ARGS[*]}
