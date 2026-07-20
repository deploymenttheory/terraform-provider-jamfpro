#!/bin/bash
# Initial cleanup
(cd ./setup && bash ./cleanup.sh)
# Scaffolding
(cd ./setup && bash ./scaffolding.sh)
# Run tests
terraform init
terraform test -parallelism=1
# Post run cleanup
(cd ./setup && bash ./cleanup.sh)
