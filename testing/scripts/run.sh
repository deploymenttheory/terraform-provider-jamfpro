#!/bin/bash

TARGETS="$1"

if [[ "$TARGETS" == *"jamfpro_static_computer_group"* ]]; then
    ./scaffolding.sh
fi

terraform fmt
terraform validate
terraform test