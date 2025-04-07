import os
import jamfpy
from dotenv import load_dotenv
from bs4 import BeautifulSoup
load_dotenv()


logger = jamfpy.get_logger(name="cleanup", level=20)
TENTANT_FQDN = "https://lbgsandbox.jamfcloud.com"

CLIENT_ID = os.environ.get("CLIENT_ID")
CLIENT_SEC = os.environ.get("CLIENT_SEC")
TESTING_ID = os.environ.get("TESTING_ID")

if TESTING_ID =="":
    logger.error("Testing ID not set correctly")
elif TESTING_ID == "local":
    logger.warning("Testing ID set to local. If run in a pipeline, this can cause unstable behaviour for other simultaneous runs.")
else:
    logger.info(f"Cleanup with testing id {TESTING_ID}")


instance = jamfpy.Tenant(
    fqdn=TENTANT_FQDN,
    auth_method="oauth2",
    client_id=CLIENT_ID,
    client_secret=CLIENT_SEC,
    token_exp_threshold_mins=1
)

def parse_error_message(html_content):
    soup = BeautifulSoup(html_content, 'html.parser')
    p_tags = soup.find_all('p')

    if len(p_tags) > 1:
        target_paragraph = p_tags[1]
        return target_paragraph.get_text()
    else:
        print("Issue parsing error response.")



def testing_ids_from_resources(resources):
    resource_ids = []
    for resource in resources:
        name = str(resource["name"])
        if name.startswith(f"tf-testing-{TESTING_ID}"):

            resource_id = resource["id"]
            resource_ids.append(resource_id)
    return resource_ids

def purge_classic_test_resources(resource_instance, resource_type_string):
    # Escape characters for underlining
    print(f"\n\033[4mPurging {resource_type_string}...\033[0m")


    resp = resource_instance.get_all()
    resp.raise_for_status()
    resources = resp.json()[resource_type_string]
    resource_ids = testing_ids_from_resources(resources)
    for id in resource_ids:
        del_resp = resource_instance.delete_by_id(id)
        if del_resp.ok:
            logger.info(f"Successfully DELETED {resource_type_string} id:{id}")
        else:
            error_response = parse_error_message(del_resp.text)
            logger.warning(f"FAILED to DELETE {resource_type_string} id:{id}\n Reason: {error_response}")


# ============================================================================ #
# Add resources to be deleted below

purge_classic_test_resources(instance.classic.scripts, "scripts")
purge_classic_test_resources(instance.classic.buildings, "buildings")
purge_classic_test_resources(instance.classic.computer_extension_attributes, "computer_extension_attributes")
purge_classic_test_resources(instance.classic.categories, "categories")
purge_classic_test_resources(instance.classic.computer_groups, "computer_groups")
purge_classic_test_resources(instance.classic.sites, "sites")
purge_classic_test_resources(instance.classic.computers, "computers")
purge_classic_test_resources(instance.classic.departments, "departments")

