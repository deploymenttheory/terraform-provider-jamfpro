import os
import jamfpy
from dotenv import load_dotenv
load_dotenv()

logger = jamfpy.get_logger(name="cleanup", level=20)
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
        name = str(resource["name"])
        if name.startswith("tf-testing"):
            resource_id = resource["id"]
            resource_ids.append(resource_id)
    return resource_ids

def purge_classic_test_resources(resource_instance, resource_type_string):
    print(f"\n######### Purging {resource_type_string} #########")

    resp = resource_instance.get_all()
    resp.raise_for_status()
    resources = resp.json()[resource_type_string]
    resource_ids = testing_ids_from_resources(resources)
    for id in resource_ids:
        del_resp = resource_instance.delete_by_id(id)
        if del_resp.ok:
            logger.info(f"Sucessfully deleted {resource_type_string} id:{id}")
        else:
            logger.warning(f"FAILED to delete {resource_type_string} id:{id}")


# ============================================================================ #
# Add resources to be deleted below

purge_classic_test_resources(instance.classic.scripts, "scripts")
purge_classic_test_resources(instance.classic.buildings, "buildings")
purge_classic_test_resources(instance.classic.computer_extension_attributes, "computer_extension_attributes")
# TODO: Add categories to jamfpy and here