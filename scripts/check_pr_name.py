#!/usr/bin/env python3

import argparse
import sys

VALID_PREFIXES = [
    "feat:",
    "fix:",
    "chore:"
]


def check_pr_name(pr_name: str):
    """
    Process the input string argument.
    
    Args:
        input_string (str): The input string to process
        
    Returns:
        None
    """
    if any(pr_name.lower().startswith(i) for i in VALID_PREFIXES):
        return
    
    print(f"PR Name has invalid prefix: {pr_name}, should be one of: {VALID_PREFIXES}")
    sys.exit(1)


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
