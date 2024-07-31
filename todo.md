To Do:

- Standardise "Criteria" schema blocks across all three advanced searches and computer groups, user groups.
- Decide if we're scoping groups to accounts or accounts to groups. Currently allowed scoping accounts to groups in accounts.
- Review all schema validation funcs and standardise if possible. 
- (SDK) Created shared struct for LDAPServer across accounts/accountgroup
- Review Computer Inventory Collection Schema.
- Computer Prestage Enrollments entire thing.
- Self Service Icons and Categories in Policies and Configuration Profiles.
- Restricted Software review if Scope can be aligned to Policy & MacOsConfigProfiles scope.
- Refactor UserGroups to mirror Computer groups logic
- Standardise construction logic from TypeList and TypeMap across all endpoints
- Standardise stating logic across all endpoints.
- Move getStringSliceFromSet function out of accountgroups and into shared package.
- Amend account privs for Jamf Pro 11.6+ (Removal of casper admin keys?)
- Adjust Account/Account Group privileges to be pulled from an automatically updated json file

Known Issues:
1. Declarative resource redeployment fails if: 
    1. Deleted resource has Site set and...
    2. Terraform Apply is run very soon (sub 30s) after resource is deleted by means other than Terraform. Jamf Pro appears to have staged deletion process in which the resource still exists but minus it's site for short period of time, confusing terraform.