#!/bin/bash
RUNID=""
FORCE=false

while getopts "r:" opt; do
  case $opt in
    r)
      RUNID=$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
  esac
done

printf "\n##### Beginning scaffolding #####\n"
printf "running computer_groups_scaffolding.py...\n"


if [ -n "$RUNID" ]; then
  python3 ../jamfpy/static_computer_groups_scaffolding.py -r $RUNID
else
  python3 ../jamfpy/static_computer_groups_scaffolding.py
fi
