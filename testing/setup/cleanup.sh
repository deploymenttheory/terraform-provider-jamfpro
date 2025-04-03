#!/bin/bash
printf "\n##### BEGINNING CLEANUP #####\n"
printf "Purging data source files in /data_sources...\n"
rm -rf ../data_sources
printf "purging resources in jamfpro...\n"
python3 ../jamfpy/clean_up.py