#!/bin/bash

if [ -z "$DOCKER_COMPOSE_BASE_COMMAND" ]; then
    echo "DOCKER_COMPOSE_BASE_COMMAND must be set"
    exit 1
fi

if [ "$#" -lt "1" ]; then
    echo "Must supply a service name as the first argument"
    exit 1
fi

# Convenience for typing it out
DC=$DOCKER_COMPOSE_BASE_COMMAND

# Get image used by service in format 'REPOSITORY:TAG' i.e. traefik:v2.6
IMAGE=$($DC images $1 | tail -n +2 | tail -n 1 | awk '{ printf("%s:%s", $2, $3) }')

case "$IMAGE" in
# Use bash wildcards to match patterns
# i.e. postgres:* matches any version of postgres image
    postgres:*)
        echo "/bin/bash"
        ;;
    dpage/pgadmin4*)
        echo "/bin/sh"
        ;;
    mmbrnr_dev_*) # Anything built by docker compose
        echo "/bin/bash"
        ;;
    traefik:*|\
    *) # Default
        # Take a punt this is pretty much always there
        echo "/bin/sh"
        ;;
esac
