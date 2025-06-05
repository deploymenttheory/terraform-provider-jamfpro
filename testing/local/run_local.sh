#!/bin/bash
# Initial cleanup
bash ./cleanup.sh
# Scaffolding
(cd ../action_scripts && python3 ./scaffolding_static_group_computers.py)
# Run tests
terraform init
terraform test
# # Post run cleanup
bash ./cleanup.sh
