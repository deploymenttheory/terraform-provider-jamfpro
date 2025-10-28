'''
This script interacts with the GitHub API to determine changed files within a specific pull request.
It then filters these files to identify 'resource targets'. A file is considered to contain
a resource target if its path includes both "internal" and "resources" segments.
The script extracts the name of the resource (assumed to be the directory name immediately
following "internal/resources/") from such paths.

All unique, extracted resource target names are then written to a file named 'targets.txt',
comma-separated. If no resource targets are found in the pull request's changed files,
the script prints a message and exits with an error code.

Command-line arguments:
  --repo-owner: The owner of the GitHub repository.
  --repo-name: The name of the GitHub repository.
  --pr-number: The number of the pull request to inspect.
  --github-token: A GitHub personal access token with permissions to read repository data.

Example Usage:

Input:
  Command:
  `python scripts/get_targets_to_file.py --repo-owner example-owner --repo-name example-repo --pr-number 123 --github-token <your_token>`

  Assuming PR #123 in `example-owner/example-repo` has the following changed files:
    - `internal/resources/user/main.go`
    - `internal/resources/group/resource.go`
    - `docs/users.md`
    - `README.md`

Output (`targets.txt`):
  `user,group`
'''

import sys
import argparse
import requests


FILEPATH_KEY = "filename"

def get_diff(owner, repo, token, pr_number: str):
    resp = requests.request(
        method="GET",
        url=f"https://api.github.com/repos/{owner}/{repo}/pulls/{pr_number}/files",
        headers={
            "Authorization": token
        }
    )
    resp.raise_for_status()
    json_resp = resp.json()
    return json_resp

def get_diff_path(response: list[dict]):
    files = []
    for i in response:
        for k, v in i.items():
            if k == FILEPATH_KEY:
                files.append(v)
    return files

def extract_resource_from_path(path: str):
    path_split = path.split("/")
    if all(i in path_split for i in ["internal", "services"]):
        return True, path_split[2]
    return False, None

def save_targets_to_file(targets: list):
    with open("targets.txt", "w", encoding="utf-8") as f:
        f.write(",".join(targets))

def main():
    parser = argparse.ArgumentParser(description="Get PR diff and extract resource targets.")
    parser.add_argument("--repo-owner", required=True, help="Repository owner.")
    parser.add_argument("--repo-name", required=True, help="Repository name.")
    parser.add_argument("--pr-number", required=True, help="Pull request number.")
    parser.add_argument("--github-token", required=True, help="GitHub token.")
    args = parser.parse_args()

    diff_info = get_diff(
        args.repo_owner,
        args.repo_name,
        args.github_token,
        args.pr_number
    )

    filepaths = get_diff_path(diff_info)

    targets = []
    for f in filepaths:
        found, res = extract_resource_from_path(f)
        if found:
            targets.append("jamfpro_" + res)

    targets = list(set(targets))

    if not targets:
        print("no targets found")
        sys.exit(1)

    save_targets_to_file(targets)

if __name__ == "__main__":
    main()
