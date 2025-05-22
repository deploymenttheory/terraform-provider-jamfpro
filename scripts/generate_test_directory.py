"""
Generates Terraform test files (.tftest.hcl) for specified resources.

This script takes a comma-separated list of resource names or the keyword 'all'
as a command-line argument. For each resource, it creates a corresponding
directory under 'testing/tests/' and generates a '<resource_name>.tftest.hcl'
file within that directory.

The content of the generated HCL file is based on a predefined template,
pointing to a payload directory expected at 'testing/payloads/<resource_name>'.

If 'all' is specified, the script discovers available resources by listing
subdirectories within 'testing/payloads/'.

Arguments:
    resources_data (str): A comma-separated string of resource names (e.g.,
                          "resource1,resource2") or the keyword 'all' to
                          process all resources found in 'testing/payloads/'.

Example:
    Input:
    `python scripts/generate_test_directory.py "foo,bar"`

    This command will:
    1. Create the directory `testing/tests/foo/`.
    2. Create the file `testing/tests/foo/foo.tftest.hcl` with content:
       ```hcl
       run "apply_foo" {
         command = apply

         module {
           source = "testing/payloads/foo"
         }
       }
       ```
    3. Create the directory `testing/tests/bar/`.
    4. Create the file `testing/tests/bar/bar.tftest.hcl` with content:
       ```hcl
       run "apply_bar" {
         command = apply

         module {
           source = "testing/payloads/bar"
         }
       }
       ```

    Input:
    `python scripts/generate_test_directory.py all`

    Assuming `testing/payloads/` contains directories `resA` and `resB`:
    1. Create the directory `testing/tests/resA/`.
    2. Create the file `testing/tests/resA/resA.tftest.hcl` (content similar to above).
    3. Create the directory `testing/tests/resB/`.
    4. Create the file `testing/tests/resB/resB.tftest.hcl` (content similar to above).

    If `testing/payloads/` is missing or empty when 'all' is used, a warning
    will be printed, and the script may exit or generate no files.
"""

import os
import sys
import argparse

TEST_BLOCK = """
run "apply_{resource_type}" {
  command = apply

  module {
    source = "{payload_dir}"
  }
}
"""

def generate_targetted_test_files(resources):
    root_dir = "testing/tests"
    os.makedirs(root_dir, exist_ok=True)

    for r in resources:
        resource_dir = os.path.join(root_dir, r)
        os.makedirs(resource_dir, exist_ok=True) # Ensure per-resource directory exists

        test_block_content = TEST_BLOCK.format(
            resource_type=r,
            payload_dir=f"testing/payloads/{r}" # Assuming payloads are structured this way
        )

        # Corrected file path to be inside the resource-specific directory
        test_file_path = os.path.join(resource_dir, f"{r}.tftest.hcl")
        with open(test_file_path, "w") as f:
            f.write(test_block_content)


def get_all_available_test_files():
    payloads_dir = "testing/payloads"
    available_resources = []

    if not os.path.isdir(payloads_dir):
        print(f"Warning: Payload directory '{payloads_dir}' not found. Cannot determine all available tests.")
        sys.exit(1)

    res_folder_list = os.listdir(payloads_dir)
    if not res_folder_list:
        print("no resource folders")
        sys.exit(1)

    for item in res_folder_list:
        available_resources.append(item)
    
    if not available_resources:
        print(f"Warning: No subdirectories found in '{payloads_dir}'. No tests will be generated for 'all'.")

    return available_resources


def main():
    parser = argparse.ArgumentParser(description="Generate targeted test files based on a resources string.")
    parser.add_argument("resources_data", help="Target resources (comma-separated string or 'all').")
    args = parser.parse_args()

    input_str: str = args.resources_data

    if input_str == "all":
        targets = get_all_available_test_files()
    else:
        targets = input_str.strip().split(",")

    generate_targetted_test_files(targets)

if __name__ == "__main__":
    main()



