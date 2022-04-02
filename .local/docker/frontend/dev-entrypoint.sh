#!/bin/bash

if [ "$#" -gt "0" ]; then
    $@
    exit "$?"
fi

npm run pre-dev
npm run dev