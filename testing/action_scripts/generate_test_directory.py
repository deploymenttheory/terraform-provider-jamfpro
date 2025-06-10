"""
Generates Terraform test files (.tftest.hcl) for specified resources.

This script takes a comma-separated list of resource names or the keyword 'all'
as a command-line argument. For each resource, it creates a corresponding
directory under 'terraform-provider-jamfpro/testing/tests/' and generates a '<resource_name>.tftest.hcl'
file within that directory.

The content of the generated HCL file is based on a predefined template.
The module source within the HCL will point to a payload directory relative
to the HCL file itself (e.g., '../../payloads/<resource_name>').

If 'all' is specified, the script discovers available resources by listing
subdirectories within 'terraform-provider-jamfpro/testing/payloads/'.

Arguments:
    resources_data (str): A comma-separated string of resource names (e.g.,
                          "resource1,resource2") or the keyword 'all' to
                          process all resources found in 'terraform-provider-jamfpro/testing/payloads/'.

Example:
    Assuming the script is run from the workspace root, and 'terraform-provider-jamfpro' is a subdirectory.
    Input:
    `python terraform-provider-jamfpro/scripts/generate_test_directory.py "foo,bar"`

    This command will:
    1. Create the directory `terraform-provider-jamfpro/testing/tests/foo/`.
    2. Create the file `terraform-provider-jamfpro/testing/tests/foo/foo.tftest.hcl` with content:
       ```hcl
       run "apply_foo" {{
         command = apply

         module {{
           source = "../../payloads/foo"
         }}
       }}
       ```
    3. Create the directory `terraform-provider-jamfpro/testing/tests/bar/`.
    4. Create the file `terraform-provider-jamfpro/testing/tests/bar/bar.tftest.hcl` with content:
       ```hcl
       run "apply_bar" {{
         command = apply

         module {{
           source = "../../payloads/bar"
         }}
       }}
       ```

    Input:
    `python terraform-provider-jamfpro/scripts/generate_test_directory.py all`

    Assuming `terraform-provider-jamfpro/testing/payloads/` contains subdirectories `resA` and `resB`:
    1. Create the directory `terraform-provider-jamfpro/testing/tests/resA/`.
    2. Create the file `terraform-provider-jamfpro/testing/tests/resA/resA.tftest.hcl` (content similar to above).
    3. Create the directory `terraform-provider-jamfpro/testing/tests/resB/`.
    4. Create the file `terraform-provider-jamfpro/testing/tests/resB/resB.tftest.hcl` (content similar to above).

    If `terraform-provider-jamfpro/testing/payloads/` is missing or empty when 'all' is used, a warning
    will be printed, and the script may exit or generate no files.
"""

import os
import sys
import argparse

TEST_BLOCK = """
run "apply_{resource_type}" {{
  command = apply

  module {{
    source = "{payload_dir}"
  }}
}}
"""

def generate_targetted_test_files(resources):
    root_dir = "testing/"
    if not os.path.exists(root_dir):
        os.makedirs(root_dir, exist_ok=True)
        print(f"Created directory: {root_dir}")

    if not resources:
        return

    for r in resources:
        payload_dir_for_hcl = f"./payloads/{r}"

        test_block_content = TEST_BLOCK.format(
            resource_type=r,
            payload_dir=payload_dir_for_hcl
        )

        test_file_path = os.path.join(root_dir, f"{r}.tftest.hcl")
        with open(test_file_path, "w") as f:
            f.write(test_block_content)
        print(f"Created test file: {test_file_path}")


def get_all_available_test_files():
    payloads_dir = "testing/payloads/"
    available_resources = []

    if not os.path.isdir(payloads_dir):
        print(f"Warning: Payload directory '{payloads_dir}' not found. Cannot determine all available tests.")
        sys.exit(1)

    res_folder_list = os.listdir(payloads_dir)
    if not res_folder_list:
        print("no resource folders")
        sys.exit(1)

    for item in res_folder_list:
        item_path = os.path.join(payloads_dir, item)
        if os.path.isdir(item_path):
            available_resources.append(item)
    
    if not available_resources:
        print(f"Warning: No subdirectories found in '{payloads_dir}'. No tests will be generated for 'all'.")

    return available_resources


def main():
    parser = argparse.ArgumentParser(description="Generate targeted test files based on a resources string.")
    parser.add_argument("resources_data", help="Target resources (comma-separated string or 'all').")
    args = parser.parse_args()

    input_str: str = args.resources_data

    if not input_str:
        sys.exit(1)

    targets = []
    if input_str.lower() == "all":
        targets = get_all_available_test_files()

    else:
        stripped_input = input_str.strip()
        

        if stripped_input:
            requested_targets = [t.strip() for t in stripped_input.split(",")]
            payloads_dir = "testing/payloads/"
            existing_targets = []
            missing_targets = []
            for t in requested_targets:
                if os.path.isdir(os.path.join(payloads_dir, t)):
                    existing_targets.append(t)
                else:
                    missing_targets.append(t)
            if missing_targets:
                print(f"Warning: The following resources do not exist in '{payloads_dir}' and will be skipped: {', '.join(missing_targets)}")
            targets = existing_targets

    if not targets:
        # This print might be useful to keep, as it indicates no files will be generated.
        # If you want it removed, I can do that in a follow-up.
        print("DEBUG: No targets specified or found. No test files will be generated.")

    generate_targetted_test_files(targets)

if __name__ == "__main__":
    main()



