#!/bin/bash
bash ./cleanup.sh
echo "Beginning scaffolding"
echo "running computer_groups_scaffolding.py..."
python3 ../jamfpy/computer_groups_scaffolding.py
