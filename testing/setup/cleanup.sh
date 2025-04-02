#!/bin/bash
echo "BEGINNING CLEANUP"
echo "Purging data source files in /data_sources..."
rm ../data_sources/*
echo "purging resources in jamfpro..."
python3 ../jamfpy/clean_up.py