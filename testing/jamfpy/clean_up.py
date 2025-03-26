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

print("BEGIN PURGE")
print("purging scripts")
resp = instance.classic.scripts.get_all()
resp.raise_for_status()
resources = resp.json()["scripts"]
resource_ids = testing_ids_from_resources(resources)
for id in resource_ids:
    del_resp = instance.classic.scripts.delete_by_id(id)
    print(del_resp.text)

print("purging computer extension attributes")
resp = instance.classic.computer_extension_attributes.get_all()
resp.raise_for_status()
computer_extension_attributes = resp.json()["computer_extension_attributes"]
resource_ids = testing_ids_from_resources(computer_extension_attributes)
for id in resource_ids:
    del_resp = instance.classic.computer_extension_attributes.delete_by_id(id)
    print(del_resp.text)

print("purging buildings")
resp = instance.classic.buildings.get_all()
resp.raise_for_status()
buildings = resp.json()["buildings"]
resource_ids = testing_ids_from_resources(buildings)
for id in resource_ids:
    del_resp = instance.classic.buildings.delete_by_id(id)
    print(del_resp.text)
