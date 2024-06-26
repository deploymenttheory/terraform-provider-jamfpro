TODO:

1. Standardise "Criteria" schema blocks across all three advanced searches
2. Decide if we're scoping groups to accounts or accounts to groups. Currently allowed scoping accounts to groups in accounts.
3. Review all schema validation funcs and standardise if possible. 



Known Issues:
1. Declarative resource redeployment fails if: 
    1. Deleted resource has Site set
    2. Terraform Apply is run very soon (sub 30s) after resource is deleted by means other than Terraform. Jamf Pro appears to have staged deletion process in which the resource still exists but minus it's site for short period of time, confusing terraform.