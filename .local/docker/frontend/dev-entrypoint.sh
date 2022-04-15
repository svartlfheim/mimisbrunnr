#!/bin/bash

if [ "$#" -gt "0" ]; then
    $@
    exit "$?"
fi

# We need to do this here, as it requires the mounted code.
# This new version of yarn is weird, but it reliaes on the .yarn folder in the 
# frontend directory to tell it what executable version to use. Only yarn 2+ has
# the plugin command.
yarn plugin import workspace-tools
yarn plugin import typescript

yarn install
yarn workspaces foreach -pAiv run start