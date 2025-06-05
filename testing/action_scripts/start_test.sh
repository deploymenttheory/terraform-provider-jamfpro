#!/bin/bash

TARGETS="$1"
UUID="$2"

if [[ "$TARGETS" == *"jamfpro_static_computer_group"* ]] || [[ "$TARGETS" == "all" ]]; then
    echo running scaffolding
    python3 ./action_scripts/scaffolding_static_group_computers.py -r "$UUID"
fi

terraform init
terraform fmt
terraform validate
terraform test