TODO:

1. Standardise "Criteria" schema blocks across all three advanced searches and computer groups, user groups.
2. Decide if we're scoping groups to accounts or accounts to groups. Currently allowed scoping accounts to groups in accounts.
3. Review all schema validation funcs and standardise if possible. 
4. (SDK) Created shared struct for LDAPServer across accounts/accountgroup
5. Flatten and optimise Computer Extension Attributes schema
6. Review Computer Inventory Collection Schema.
7. Computer Prestage Enrollments entire thing.
8. Self Service Icons and Categories in Policies and Configuration Profiles.
9. Restricted Software review if Scope can be aligned to Policy & MacOsConfigProfiles scope.
10. Refactor UserGroups to mirror Computer groups logic
11. Standardise construction logic from TypeList and TypeMap across all endpoints
12. Standardise stating logic across all endpoints.
13. Move getStringSliceFromSet function out of accountgroups and into shared package.
14. Amend account privs for Jamf Pro 11.6+ (Removal of casper admin keys?)
15. Remove Get by name fallback in all occurances.
16. Line 306 macosconfigurationprofiles.state
17. Finish Policies

Known Issues:
1. Declarative resource redeployment fails if: 
    1. Deleted resource has Site set and...
    2. Terraform Apply is run very soon (sub 30s) after resource is deleted by means other than Terraform. Jamf Pro appears to have staged deletion process in which the resource still exists but minus it's site for short period of time, confusing terraform.