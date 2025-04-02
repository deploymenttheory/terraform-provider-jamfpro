#!/bin/bash
(cd ./setup && bash ./scaffolding.sh)
terraform init
terraform test
