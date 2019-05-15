#!/bin/bash

PACKAGE_MGR="$1"
PACKAGE_NAME="$2"

if [[ -z "$PACKAGE_MGR" ]] || [[ -z "$PACKAGE_NAME" ]]; then
    echo "Usage: $0 <package_manager> <package_name>"
    exit 1
fi

if [ "$PACKAGE_MGR" = "gem" ]; then
    if ! gem list | grep "$PACKAGE_NAME "; then
        echo "gem $PACKAGE_NAME not installed, installing"
        gem install "$PACKAGE_NAME"
    elif gem outdated | grep "$PACKAGE_NAME "; then
        echo "gem $PACKAGE_NAME out of date, updating"
        gem update "$PACKAGE_NAME"
    fi
elif [ "$PACKAGE_MGR" = "brew" ]; then
    if ! brew list | grep "^$PACKAGE_NAME\$"; then
        echo "brew package $PACKAGE_NAME not installed, installing"
        brew install "$PACKAGE_NAME"
    elif brew outdated | grep "^$PACKAGE_NAME\$"; then
        brew upgrade "$PACKAGE_NAME"
    fi
fi

