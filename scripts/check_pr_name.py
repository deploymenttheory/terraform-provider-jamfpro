"""simple script to make sure the PR names are formed how we like"""

#!/usr/bin/env python3

import argparse
import sys
import re

VALID_PREFIXES_LOWER = [
    "feat",
    "fix",
    "chore",
    "docs",
    "style",
    "refactor",
    "perf",
    "test",
    "build",
    "ci",
    "revert"
]

class InvalidPrName(Exception):
    "Error for a bad PR name"

def check_pr_name(pr_name: str):
    """
    Process the input string argument.

    Args:
        input_string (str): The input string to process

    Returns:
        None
    """

    for i in VALID_PREFIXES_LOWER:
        matched = re.match(fr'^{i}:', pr_name)
        if matched:
            print(f"[DEBUG] successfully matched {matched.group(0)}")
            return

        matched = re.match(fr'^{i}\b', pr_name)
        if matched:
            raise InvalidPrName(
                f"[DEBUG] matched {matched.group(0)} but found no colon. Should be: '{pr_name}:'"
            )


    raise InvalidPrName(
        f"PR Name has invalid prefix: {pr_name}, should be one of: {VALID_PREFIXES_LOWER}"
        )


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
