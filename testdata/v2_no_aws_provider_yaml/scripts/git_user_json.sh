#!/bin/sh

echo "{\"name\": \"$(git --no-pager config --global --get user.name)\"}" || echo "{ \"user\": \"unknown\" }"