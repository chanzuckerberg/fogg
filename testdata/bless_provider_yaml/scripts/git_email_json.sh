#!/bin/sh

echo "{\"email\": \"$(git --no-pager config --global --get user.email)\"}" || echo "{ \"email\": \"unknown\" }"