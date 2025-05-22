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

print("DEBUG: Script starting...")

TEST_BLOCK = """
run "apply_{resource_type}" {
  command = apply

  module {
    source = "{payload_dir}"
  }
}
"""
print(f"DEBUG: TEST_BLOCK template defined as: \\n{TEST_BLOCK}")

def generate_targetted_test_files(resources):
    print(f"DEBUG: generate_targetted_test_files called with resources: {resources}")
    root_dir = "testing/tests"
    print(f"DEBUG: Root directory for tests set to: {root_dir}")
    
    print(f"DEBUG: Ensuring root directory '{root_dir}' exists...")
    os.makedirs(root_dir, exist_ok=True)
    print(f"DEBUG: Root directory '{root_dir}' ensured.")

    if not resources:
        print("DEBUG: No resources provided to generate_targetted_test_files. Exiting function.")
        return

    for r in resources:
        print(f"DEBUG: Processing resource: {r}")
        resource_dir = os.path.join(root_dir, r)
        print(f"DEBUG: Resource directory path: {resource_dir}")

        print(f"DEBUG: Ensuring resource directory '{resource_dir}' exists...")
        os.makedirs(resource_dir, exist_ok=True) # Ensure per-resource directory exists
        print(f"DEBUG: Resource directory '{resource_dir}' ensured.")

        payload_dir_rel = f"testing/payloads/{r}"
        print(f"DEBUG: Relative payload directory for resource '{r}': {payload_dir_rel}")

        test_block_content = TEST_BLOCK.format(
            resource_type=r,
            payload_dir=payload_dir_rel # Assuming payloads are structured this way
        )
        print(f"DEBUG: Generated test block content for resource '{r}':\\n{test_block_content}")

        # Corrected file path to be inside the resource-specific directory
        test_file_path = os.path.join(resource_dir, f"{r}.tftest.hcl")
        print(f"DEBUG: Test file path for resource '{r}': {test_file_path}")

        print(f"DEBUG: Writing test block content to {test_file_path}...")
        with open(test_file_path, "w") as f:
            f.write(test_block_content)
        print(f"DEBUG: Successfully wrote to {test_file_path}")
    
    print(f"DEBUG: Finished processing all resources in generate_targetted_test_files.")


def get_all_available_test_files():
    print("DEBUG: get_all_available_test_files called.")
    payloads_dir = "testing/payloads/"
    print(f"DEBUG: Payloads directory set to: {payloads_dir}")
    available_resources = []
    print(f"DEBUG: Initialized available_resources: {available_resources}")
    print(os.listdir())
    print(f"DEBUG: Checking if payload directory '{payloads_dir}' exists and is a directory...")
    if not os.path.isdir(payloads_dir):
        print(f"DEBUG: Payload directory '{payloads_dir}' not found or not a directory.")
        print(f"Warning: Payload directory '{payloads_dir}' not found. Cannot determine all available tests.")
        sys.exit(1)
    print(f"DEBUG: Payload directory '{payloads_dir}' found.")

    print(f"DEBUG: Listing contents of '{payloads_dir}'...")
    res_folder_list = os.listdir(payloads_dir)
    print(f"DEBUG: Contents of '{payloads_dir}': {res_folder_list}")

    if not res_folder_list:
        print("DEBUG: No items found in payload directory.")
        print("no resource folders")
        sys.exit(1)
    print(f"DEBUG: Found items in payload directory: {res_folder_list}")

    for item in res_folder_list:
        item_path = os.path.join(payloads_dir, item)
        print(f"DEBUG: Checking item: {item} at path: {item_path}")
        if os.path.isdir(item_path):
            print(f"DEBUG: Item '{item}' is a directory. Adding to available_resources.")
            available_resources.append(item)
        else:
            print(f"DEBUG: Item '{item}' is not a directory. Skipping.")
    
    print(f"DEBUG: Final list of available_resources (subdirectories): {available_resources}")

    if not available_resources:
        print(f"DEBUG: No subdirectories found in '{payloads_dir}' after filtering.")
        print(f"Warning: No subdirectories found in '{payloads_dir}'. No tests will be generated for 'all'.")

    print(f"DEBUG: Returning available_resources: {available_resources}")
    return available_resources


def main():
    print("DEBUG: main function called.")
    parser = argparse.ArgumentParser(description="Generate targeted test files based on a resources string.")
    print("DEBUG: ArgumentParser created.")
    parser.add_argument("resources_data", help="Target resources (comma-separated string or 'all').")
    print("DEBUG: 'resources_data' argument added to parser.")
    
    print("DEBUG: Parsing command-line arguments...")
    args = parser.parse_args()
    print(f"DEBUG: Arguments parsed: {args}")

    input_str: str = args.resources_data
    print(f"DEBUG: Input string 'resources_data': {input_str}")

    targets = []
    if input_str.lower() == "all": # Made 'all' check case-insensitive
        print("DEBUG: Input string is 'all'. Calling get_all_available_test_files().")
        targets = get_all_available_test_files()
        print(f"DEBUG: Targets received from get_all_available_test_files(): {targets}")
    else:
        print(f"DEBUG: Input string is not 'all'. Processing as comma-separated list: '{input_str}'")
        stripped_input = input_str.strip()
        print(f"DEBUG: Input string after strip(): '{stripped_input}'")
        if stripped_input: # Ensure not empty string after strip
            targets = stripped_input.split(",")
            print(f"DEBUG: Targets after split(','): {targets}")
            # Further strip individual targets if they might have spaces e.g. "res1, res2"
            targets = [t.strip() for t in targets if t.strip()]
            print(f"DEBUG: Targets after individual stripping and filtering empty strings: {targets}")
        else:
            print("DEBUG: Input string was empty or only whitespace. Targets list is empty.")
            targets = []


    if not targets:
        print("DEBUG: No targets determined. Exiting main function without generating files.")
        # Potentially add a sys.exit(0) or a message if no files being generated is an issue
        # For now, it will just call generate_targetted_test_files with an empty list.
        # generate_targetted_test_files already has a check for empty resources.
        print("DEBUG: No targets specified or found. No test files will be generated.")

    print(f"DEBUG: Calling generate_targetted_test_files with targets: {targets}")
    generate_targetted_test_files(targets)
    print("DEBUG: generate_targetted_test_files finished.")

    print("DEBUG: main function finished.")

if __name__ == "__main__":
    print(f"DEBUG: Script execution started as __main__ (name: {__name__})")
    main()
    print("DEBUG: Script execution finished.")



