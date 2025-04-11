# Provider Testing
**PLEASE NOTE at 23:59 a wipe job is performed on all testing objects. Tests run at this time will most likely fail.**
## Environment
### Running the tests

from `terraform-provider-jamfpro/testing`, please create a python virtual environment (venv) and install the requirements.txt.

```bash
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

After this, run `run_tests.sh`. This will provision and run the tests.

NOTE: running `terraform test` may fail due to missing resources. run the test script above instead.

### Symlinks

In each tesing module, there is a symlink to the testing root provider.tf. This is required for the tests to pass. 

When adding new resource tests, please be sure that a symlink is correctly produced.

This can be achieved by inputting this command in the resource directory.


` ln -s ../../provider.tf provider.tf `


The same result can be achieved by copying the file into the directories, without the benefit of only needing to alter one file for changes.

## Adding tests

Find the resource type you want to write tests for. It will be under `testing/<resource-name>`. If the resource type is not there...

1) Create a directory for the resource type and create a terraform configuration file. Write your test config in this file.
2) Symlink the root provider.tf in the directory
3) In `testing/jamfpy` add the resource to the clean up script.
4) in `tests/smoke_test.tftest.hcl` add a run block pointing to the directory made in step 1.

### Test case names

To ensure only orphaned testing resources are removed in the clean up process, all testing resources must have a prefix of `tf-testing`. For example:

```hcl
resource "jamfpro_script" "min_script" {
  name = "tf-testing-script-min"
  script_contents = "script_contents_field"
  priority = "BEFORE"
}
```



