#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

IPV4_TARGET="127.0.0.1"

sudo $SCRIPT_DIR/hosts remove host mimisbrunnr.local
sudo $SCRIPT_DIR/hosts remove host pgadmin.local

sudo $SCRIPT_DIR/hosts add $IPV4_TARGET mimisbrunnr.local
sudo $SCRIPT_DIR/hosts add $IPV4_TARGET pgadmin.local
