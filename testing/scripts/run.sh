#!/bin/bash

TARGETS="$1"
UUID="$2"

if [[ "$TARGETS" == *"jamfpro_static_computer_group"* ]]; then
    ./scaffolding.sh -r "$UUID"
fi

terraform init
terraform fmt
terraform validate
terraform test