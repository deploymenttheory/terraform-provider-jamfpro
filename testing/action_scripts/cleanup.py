import os
import sys
from optparse import OptionParser
import jamfpy
from dotenv import load_dotenv
from bs4 import BeautifulSoup
load_dotenv()


parser = OptionParser()
parser.add_option("-f", "--force", action="store_true", dest="force",
                    default=False,
                    help="Force cleanup of all tf-testing resources")
parser.add_option("-r", "--runid", dest="runid",
                    help="ID associated with the tests to clean up")
parser.add_option("-l", "--log-level", dest="loglevel",
                  default=20,
                  help="DEBUG: 10\nINFO: 20\n WARN: 30")


(options, args) = parser.parse_args()

logger = jamfpy.get_logger(name="cleanup", level=options.loglevel)

if options.force and options.runid:
    print("-f or --force overrides runid")
    sys.exit()
elif options.force:
    TESTING_ID = None
    logger.info("CLEANING ALL tf-testing* RESOURCES. This will affect any tests runs currently in progress.")
elif options.runid:
    TESTING_ID = options.runid
    logger.info(f"Cleaning tf-testing-{TESTING_ID}*")
else:
    TESTING_ID = "local"
    logger.warning(f"Cleaning tf-testing-{TESTING_ID}* If this is running in a pipeline, this script is being called incorrectly.")


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


def parse_error_message(html_content):
    soup = BeautifulSoup(html_content, 'html.parser')
    p_tags = soup.find_all('p')

    if len(p_tags) > 1:
        target_paragraph = p_tags[1]
        return target_paragraph.get_text()

    logger.warning("Issue parsing error response.")


def testing_ids_from_resources(resources):
    resource_ids = []
    if TESTING_ID:
        prefix = f"tf-testing-{TESTING_ID}"
    else:
        prefix = "tf-testing"

    for resource in resources:
        name = str(resource["name"])
        if name.startswith(prefix):
            resource_id = resource["id"]
            resource_ids.append(resource_id)
    return resource_ids


def purge_classic_test_resources(resource_instance, resource_type_string):
    # Escape characters for underlining
    logger.info(f"\033[4mPurging {resource_type_string}...\033[0m")


    resp = resource_instance.get_all()
    resp.raise_for_status()
    resources = resp.json()[resource_type_string]
    resource_ids = testing_ids_from_resources(resources)
    for res_id in resource_ids:
        del_resp = resource_instance.delete_by_id(res_id)
        if del_resp.ok:
            logger.info(f"Successfully DELETED {resource_type_string} id:{res_id}")
        else:
            error_response = parse_error_message(del_resp.text)
            logger.warning(f"FAILED to DELETE {resource_type_string} id:{res_id}\n Reason: {error_response}")


# ============================================================================ #
# Add resources to be deleted below
purge_classic_test_resources(instance.classic.scripts, "scripts")
purge_classic_test_resources(instance.classic.buildings, "buildings")
purge_classic_test_resources(instance.classic.computer_extension_attributes, "computer_extension_attributes")
purge_classic_test_resources(instance.classic.categories, "categories")
purge_classic_test_resources(instance.classic.computer_groups, "computer_groups")
purge_classic_test_resources(instance.classic.mobile_device_groups, "mobile_device_groups")
purge_classic_test_resources(instance.classic.sites, "sites")
purge_classic_test_resources(instance.classic.computers, "computers")
purge_classic_test_resources(instance.classic.departments, "departments")
