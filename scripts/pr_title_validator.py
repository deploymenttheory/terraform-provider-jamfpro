"""Makes sure pr titles are valid"""

import os

# pre reqs
class InvalidPRTitle(Exception):
    """title is invalid"""

VALID_PREFIXES = [
    "feat",
    "fix",
    "docs",
    "style",
    "refactor",
    "perf",
    "test",
    "build",
    "ci",
    "chore",
    "revert"
]

# Make sure it exists
PR_TITLE = os.environ.get("PR_TITLE")

if not PR_TITLE:
    raise InvalidPRTitle("no PR title was provided, please check the env var value under PR_TITLE")

# Tests

# Needs a colon, split here to make below easier
if ":" not in PR_TITLE:
    raise InvalidPRTitle("no colon found in title")

PR_TITLE_PREFIX = PR_TITLE.split(":", maxsplit=1) + ":"

# No upper case - release please doesn't like it.
if any(char.isupper() for char in PR_TITLE_PREFIX):
    raise InvalidPRTitle(f"title prefix contains uppercase chars: {PR_TITLE}")

if not any(i in PR_TITLE for i in VALID_PREFIXES):
    raise InvalidPRTitle("title does not contain a conventional commit token")

if "()" in PR_TITLE_PREFIX:
    raise InvalidPRTitle("cannot have empty scope")

if " " in PR_TITLE_PREFIX:
    raise InvalidPRTitle("cannot have spaces in title prefix")

if (len(PR_TITLE) - len(PR_TITLE_PREFIX)) < 5:
    raise InvalidPRTitle(f"title message must be more than 5 chars: {PR_TITLE}")


print(f"title valid: {PR_TITLE}")