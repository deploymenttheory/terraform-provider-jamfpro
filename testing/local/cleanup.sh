#!/bin/bash
RUNID=""
FORCE=false

while getopts "r:f" opt; do
  case $opt in
    r)
      RUNID=$OPTARG
      ;;
    f)
      FORCE=true
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
  esac
done

printf "\n##### BEGINNING CLEANUP #####\n"
printf "Purging data source files in /data_sources...\n"
rm -rf ../data_sources
printf "purging resources in jamfpro...\n"

if [ "$FORCE" = true ]; then
  python3 jamfpy/clean_up.py -f
elif [ -n "$RUNID" ]; then
  python3 jamfpy/clean_up.py -r $RUNID
else
  python3 jamfpy/clean_up.py
fi
