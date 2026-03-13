#!/bin/sh

git --no-pager log --pretty=format:"{ \"sha\": \"%H\" }" -1 HEAD  || echo "{ \"sha\": \"unknown\" }"