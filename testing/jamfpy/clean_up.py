import os
import jamfpy
from dotenv import load_dotenv
load_dotenv()
TENTANT_FQDN = "https://lbgsandbox.jamfcloud.com"

CLIENT_ID = os.environ.get("CLIENT_ID")
CLIENT_SEC = os.environ.get("CLIENT_SEC")

instance = jamfpy.Tenant(
    fqdn=TENTANT_FQDN,
    auth_method="oauth2",
    client_id=CLIENT_ID,
    client_secret=CLIENT_SEC,
    token_exp_threshold_mins=1
)

def testing_ids_from_resources(resources):
    resource_ids = []
    for resource in resources:
        name = resource["name"]
        prefix = name[0:10]
        if prefix == "tf-testing":
            resource_id = resource["id"]
            resource_ids.append(resource_id)
    return resource_ids

def purge_classic_test_resources(resource_instance, resource_type_string):
    print(f"######### Purging {resource_type_string} #########")

    resp = resource_instance.get_all()
    resp.raise_for_status()
    resources = resp.json()[resource_type_string]
    resource_ids = testing_ids_from_resources(resources)
    for id in resource_ids:
        del_resp = resource_instance.delete_by_id(id)
        print(del_resp.text)

# ============================================================================ #
# Add resources to be deleted below

purge_classic_test_resources(instance.classic.scripts, "scripts")
purge_classic_test_resources(instance.classic.buildings, "buildings")
purge_classic_test_resources(instance.classic.computer_extension_attributes, "computer_extension_attributes")