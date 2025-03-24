# Provider Testing

In each tesing module, there is a symlink to the testing root main.tf and terraform.tfvars. This is required for the tests to pass. 

When adding new resource tests, please be sure that a symlink is correctly produced.

This can be achieved by inputting this command in the resource directory.


` ln -s ../provider.tf provider.tf `


The same result can be achieved by copying the file into the directories, without the benefit of only needing to alter one file for changes.