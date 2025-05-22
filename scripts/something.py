import os
import sys
import requests
from dotenv import load_dotenv

REQUIRED_ENV_VARS = [
    "REPO_OWNER",
    "REPO_NAME",
    "PR_NUMBER",
    "GITHUB_TOKEN"
]

FILEPATH_KEY = "filename"

def get_env():
    load_dotenv()
    out = {}
    for k in REQUIRED_ENV_VARS:
        v = os.environ.get(k)
        if not v:
            print(f"::error::Required environment variable {k} is not set")
            sys.exit(1)
        
        out[k] = v

    return out


def get_diff(owner, repo, token, pr_number: str):
    resp = requests.request(
        method="GET",
        url=f"https://api.github.com/repos/{owner}/{repo}/pulls/{pr_number}/files",
        headers={
            "Authorization": token
        }
    )

    resp.raise_for_status()

    return resp.json()


def get_diff_path(response: list[dict]):
    files = []
    for i in response:
        for k, v in i.items():
            if k == FILEPATH_KEY:
                files.append(v)

    return files


def extract_resource_from_path(path: str):
    path_split = path.split("/")
    if all(i in path_split for i in ["internal", "resources"]):
        return True, path_split[2]
    
    return False, None


def save_targets_to_file(targets: list):
    with open("target_resources.txt", "w") as f:
        f.write(",".join(targets))        


def main():
    env = get_env()
    diff_info = get_diff(
        env["REPO_OWNER"],
        env["REPO_NAME"],
        env["GITHUB_TOKEN"],
        env["PR_NUMBER"]
    )

    filepaths = get_diff_path(diff_info)

    targets = []
    for f in filepaths:
        found, res = extract_resource_from_path(f)

        if found:
            targets.append(res)

    save_targets_to_file(targets)


if __name__ == "__main__":
    main()