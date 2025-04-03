#!/bin/bash
bash ./cleanup.sh
echo "Beginning scaffolding"
echo "running computer_groups_scaffolding.py..."
python3 ../jamfpy/static_computer_groups_scaffolding.py
