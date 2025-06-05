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

import os
import sys
import requests
import argparse

print("[DEBUG] Starting script execution.")

FILEPATH_KEY = "filename"

def get_diff(owner, repo, token, pr_number: str):
    print(f"[DEBUG] get_diff called with owner={owner}, repo={repo}, pr_number={pr_number}")
    resp = requests.request(
        method="GET",
        url=f"https://api.github.com/repos/{owner}/{repo}/pulls/{pr_number}/files",
        headers={
            "Authorization": token
        }
    )
    print(f"[DEBUG] GitHub API response status: {resp.status_code}")
    resp.raise_for_status()
    json_resp = resp.json()
    return json_resp


def get_diff_path(response: list[dict]):
    print(f"[DEBUG] get_diff_path called with response: {response}")
    files = []
    for i in response:
        print(f"[DEBUG] Inspecting diff item: {i}")
        for k, v in i.items():
            print(f"[DEBUG] Key: {k}, Value: {v}")
            if k == FILEPATH_KEY:
                print(f"[DEBUG] Found file path: {v}")
                files.append(v)
    print(f"[DEBUG] Extracted file paths: {files}")
    return files


def extract_resource_from_path(path: str):
    print(f"[DEBUG] extract_resource_from_path called with path: {path}")
    path_split = path.split("/")
    print(f"[DEBUG] Path split: {path_split}")
    if all(i in path_split for i in ["internal", "resources"]):
        print(f"[DEBUG] Path contains both 'internal' and 'resources'. Resource: {path_split[2]}")
        return True, path_split[2]
    print(f"[DEBUG] Path does not contain both 'internal' and 'resources'.")
    return False, None


def save_targets_to_file(targets: list):
    print(f"[DEBUG] save_targets_to_file called with targets: {targets}")
    with open("targets.txt", "w") as f:
        f.write(",".join(targets))        
    print(f"[DEBUG] targets.txt written with: {','.join(targets)}")


def main():
    print("[DEBUG] main() called.")
    parser = argparse.ArgumentParser(description="Get PR diff and extract resource targets.")
    parser.add_argument("--repo-owner", required=True, help="Repository owner.")
    parser.add_argument("--repo-name", required=True, help="Repository name.")
    parser.add_argument("--pr-number", required=True, help="Pull request number.")
    parser.add_argument("--github-token", required=True, help="GitHub token.")
    args = parser.parse_args()
    print(f"[DEBUG] Parsed arguments: {args}")

    diff_info = get_diff(
        args.repo_owner,
        args.repo_name,
        args.github_token,
        args.pr_number
    )
    print(f"[DEBUG] diff_info: {diff_info}")

    filepaths = get_diff_path(diff_info)
    print(f"[DEBUG] filepaths: {filepaths}")

    targets = []
    for f in filepaths:
        print(f"[DEBUG] Processing file path: {f}")
        found, res = extract_resource_from_path(f)
        print(f"[DEBUG] extract_resource_from_path result: found={found}, res={res}")
        if found:
            targets.append(res)
            print(f"[DEBUG] Appended target: {res}")

    print(f"[DEBUG] Final targets list: {targets}")
    if not targets:
        print("[DEBUG] No targets found. Exiting with error.")
        print("no targets found")
        sys.exit(1)

    save_targets_to_file(targets)
    print("[DEBUG] Script completed successfully.")


if __name__ == "__main__":
    print("[DEBUG] __main__ entrypoint.")
    main()