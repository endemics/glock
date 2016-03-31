#!/bin/sh
# try to get a lock from a remote glock locking service
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

# Max wait time: LOOP * WAIT seconds
# DEFAULT = 30 minutes
LOOP=180  # how many times do we retry?
WAIT=10   # in seconds

# return codes from glock when trying to lock
LOCK_OK=201
LOCK_KO=409

lock_it () {
  RET=`curl -s -o /dev/null -w "%{http_code}" -X PUT ${URL}${LOCK}`
  if [ $RET -eq ${LOCK_OK} ]; then
    return 0
  elif [ $RET -eq ${LOCK_KO} ]; then
    return 1
  else
    echo "Remote locking service returned an error, aborting" >&2
    exit 1
  fi
}

i=0

while [ $i -lt ${LOOP} ];do
  if lock_it; then
    echo "lock \"${LOCK}\" acquired"
    exit 0
  fi
  printf '.'
  i=`expr $i + 1`
  sleep ${WAIT}
done

exit 1
