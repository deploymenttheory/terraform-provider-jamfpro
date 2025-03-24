# Provider Testing

In each tesing module, there is a symlink to the testing root main.tf and terraform.tfvars. This is required for the tests to pass. 

When adding new resource tests, please be sure that a symlink is correctly produced.