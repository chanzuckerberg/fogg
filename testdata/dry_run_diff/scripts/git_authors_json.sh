#!/bin/sh

echo "{\"authors\": \"$(git --no-pager log --format='%an' -- . | grep -v "[bot]" | head -10)\"}" || echo "{ \"authors\": \"unknown\" }"