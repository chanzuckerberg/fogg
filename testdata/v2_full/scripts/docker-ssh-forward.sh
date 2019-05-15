#!/usr/bin/env bash
set -eo pipefail

# copied from https://github.com/uber-common/docker-ssh-agent-forward/blob/master/pinata-ssh-forward.sh

if ! ssh-add -L >/dev/null; then
  echo "Add ssh key(s) to your agent and try again."
  echo "This is needed for pulling modules from github."
  exit
fi

IMAGE_NAME=uber/ssh-agent-forward:latest
CONTAINER_NAME=pinata-sshd
VOLUME_NAME=ssh-agent
HOST_PORT=2244
AUTHORIZED_KEYS=$(ssh-add -L | base64 | tr -d '\n')
KNOWN_HOSTS_FILE=$(mktemp -t dsaf.XXX)

trap 'rm ${KNOWN_HOSTS_FILE}' EXIT

docker rm -f "${CONTAINER_NAME}" >/dev/null 2>&1 || true

docker volume create --name "${VOLUME_NAME}"

docker run \
  --name "${CONTAINER_NAME}" \
  -e AUTHORIZED_KEYS="${AUTHORIZED_KEYS}" \
  -v ${VOLUME_NAME}:/ssh-agent \
  -d \
  -p "${HOST_PORT}:22" \
  "${IMAGE_NAME}" >/dev/null \
;

if [ "${DOCKER_HOST}" ]; then
  HOST_IP=$(echo "$DOCKER_HOST" | awk -F '//' '{print $2}' | awk -F ':' '{print $1}')
else
  HOST_IP=127.0.0.1
fi

# FIXME Find a way to get rid of this additional 1s wait
sleep 1
while ! nc -z -w5 ${HOST_IP} ${HOST_PORT}; do sleep 0.1; done

ssh-keyscan -p "${HOST_PORT}" "${HOST_IP}" >"${KNOWN_HOSTS_FILE}" 2>/dev/null

# show the keys that are being forwarded
ssh \
  -A \
  -o "UserKnownHostsFile=${KNOWN_HOSTS_FILE}" \
  -p "${HOST_PORT}" \
  -S none \
  "root@${HOST_IP}" \
  ssh-add -l

# keep the agent running
ssh \
  -A \
  -f \
  -o "UserKnownHostsFile=${KNOWN_HOSTS_FILE}" \
  -p "${HOST_PORT}" \
  -S none \
  "root@${HOST_IP}" \
  /ssh-entrypoint.sh

echo 'Agent forwarding successfully started.'