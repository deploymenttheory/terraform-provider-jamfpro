#!/usr/bin/env python3

import argparse
import sys
import re

VALID_PREFIXES_LOWER = [
    "feat",
    "fix",
    "chore"
]


def check_pr_name(pr_name: str):
    """
    Process the input string argument.
    
    Args:
        input_string (str): The input string to process
        
    Returns:
        None
    """
    if any(pr_name.lower().startswith(i) for i in VALID_PREFIXES_LOWER):
        return
    
    
    for i in VALID_PREFIXES_LOWER:
        matched = re.match(fr'^{i}:')
        if matched:
            print(f"[DEBUG] successfully matched {matched.group(1)}")
            return
    
    raise Exception(f"PR Name has invalid prefix: {pr_name}, should be one of: {VALID_PREFIXES_LOWER}")


def main():
    """Main entry point for the script."""
    parser = argparse.ArgumentParser(description='Process a single string argument.')
    parser.add_argument('pr_name', type=str, help='Input string to process')
    
    pr_name = parser.parse_args().pr_name

    if not pr_name:
        print("Error: PR name is empty")
        sys.exit(1)

    check_pr_name(pr_name)

if __name__ == "__main__":
    main()
