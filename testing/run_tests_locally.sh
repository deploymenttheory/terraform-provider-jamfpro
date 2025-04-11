#!/bin/bash
# Initial cleanup
(cd ./setup && bash ./cleanup.sh)
# Scaffolding
(cd ./setup && bash ./scaffolding.sh)
# Run tests
terraform init
terraform test
# Post run cleanup
(cd ./setup && bash ./cleanup.sh)
