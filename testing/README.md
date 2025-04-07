# Provider Testing

## Environment
In each tesing module, there is a symlink to the testing root provider.tf. This is required for the tests to pass. 

When adding new resource tests, please be sure that a symlink is correctly produced.

This can be achieved by inputting this command in the resource directory.


` ln -s ../../provider.tf provider.tf `


The same result can be achieved by copying the file into the directories, without the benefit of only needing to alter one file for changes.

To begin the tests, from `/testing` run 

`terraform init`

`terraform test`

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



