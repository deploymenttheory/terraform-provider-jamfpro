#!/bin/bash

TARGETS="$1"
UUID="$2"

if [[ "$TARGETS" == *"jamfpro_static_computer_group"* ]]; then
    echo running scaffolding
    ./scripts/computer_groups_scaffolding.py -r "$UUID"
fi

terraform init
terraform fmt
terraform validate
terraform test