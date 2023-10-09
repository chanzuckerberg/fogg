#!/bin/bash

# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

# I would have written this directly in the Makefile, but that was difficult.

CMD="$1"

TMP=$(mktemp)
TMP2=$(mktemp)
terraform-docs md . >"$TMP"
sed '/^<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->$/,/<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->/{//!d;}' README.md | sed "/^<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->$/r $TMP" >$TMP2

case "$CMD" in
update)
    mv $TMP2 README.md
    ;;
check)
    diff $TMP2 README.md >/dev/null
    ;;
esac

exit $?