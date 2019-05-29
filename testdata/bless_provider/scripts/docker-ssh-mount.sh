#!/usr/bin/env bash

# copied from https://raw.githubusercontent.com/uber-common/docker-ssh-agent-forward/master/pinata-ssh-mount.sh

echo "-v ssh-agent:/ssh-agent -e SSH_AUTH_SOCK=/ssh-agent/ssh-agent.sock"
