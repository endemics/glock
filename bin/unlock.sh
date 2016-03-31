#!/bin/sh
# try to remove a lock from a remote glock locking service
set -e

usage () {
  echo "Usage: $0 <glock url> <lockname>"
}

if [ $# -ne 2 ]; then
  usage
  exit 1
fi

URL=$1   # The whole URL, including proto and path
LOCK=$2

# return codes from glock when trying to unlock
UNLOCK_OK=200

unlock_it () {
  RET=`curl -s -o /dev/null -w "%{http_code}" -X DELETE ${URL}${LOCK}`
  if [ $RET -eq ${UNLOCK_OK} ]; then
    echo "lock \"${LOCK}\" removed"
    exit 0
  else
    echo "Remote locking service returned an error, aborting" >&2
    exit 1
  fi
}

unlock_it
