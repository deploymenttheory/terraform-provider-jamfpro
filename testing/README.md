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

## Test case names

To ensure only orphaned testing resources are removed in the clean up process, all testing resources must have a prefix of `tf-testing`. For example:

```hcl
resource "jamfpro_script" "min_script" {
  name = "tf-testing-script-min"
  script_contents = "script_contents_field"
  priority = "BEFORE"
}
```

