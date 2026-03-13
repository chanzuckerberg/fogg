#!/bin/sh

echo "{\"branch\": \"$(git --no-pager rev-parse --abbrev-ref HEAD)\"}" || echo "{ \"branch\": \"unknown\" }"