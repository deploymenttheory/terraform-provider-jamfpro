import os
import random
import uuid
import json
from optparse import OptionParser
from pathlib import Path
import jamfpy
from dotenv import load_dotenv
load_dotenv()
print("SCAFFOLDING")
logger = jamfpy.new_logger(name="site_mobile_device_setup", level=20)

parser = OptionParser()
parser.add_option("-r", "--runid", dest="runid",
                    help="Create scaffolding objects with given ID")

(options, args) = parser.parse_args()

if options.runid:
    TESTING_ID = options.runid
    logger.info(f"Creating scaffolding objects with tf-testing-{TESTING_ID}*")
else:
    TESTING_ID = "local"
    logger.warning(f"Creating scaffolding objects with tf-testing-{TESTING_ID}* If run in a pipeline, this script is being called incorrectly.")


TENTANT_FQDN = "https://lbgsandbox.jamfcloud.com"
CLIENT_ID = os.environ.get("CLIENT_ID")
CLIENT_SEC = os.environ.get("CLIENT_SEC")
RANDOM_NUMBER = random.randint(0,9999)
DEVICE_COUNT = 10
SITE_NAME = f"tf-testing-{TESTING_ID}-site-{RANDOM_NUMBER}"


instance = jamfpy.Tenant(
    fqdn=TENTANT_FQDN,
    auth_method="oauth2",
    client_id=CLIENT_ID,
    client_secret=CLIENT_SEC,
    token_exp_threshold_mins=1
)


def create_mobile_device_config(device_name, site_id, site_name):
    return f"""
<mobile_device>
    <general>
        <name>{device_name}</name>
        <serial_number>{uuid.uuid4()}</serial_number>
        <udid>{uuid.uuid4()}</udid>
        <asset_tag/>
        <site>
            <id>{site_id}</id>
            <name>{site_name}</name>
        </site>
    </general>
    <location>
        <username/>
        <realname/>
        <real_name/>
        <email_address/>
        <position/>
        <phone/>
        <phone_number/>
        <department/>
        <building/>
        <room/>
    </location>
    <purchasing>
        <is_purchased>true</is_purchased>
        <is_leased>false</is_leased>
        <po_number/>
        <vendor/>
        <purchase_price/>
        <purchasing_account/>
        <po_date/>
        <po_date_epoch>0</po_date_epoch>
        <po_date_utc/>
        <warranty_expires/>
        <warranty_expires_epoch>0</warranty_expires_epoch>
        <warranty_expires_utc/>
        <lease_expires/>
        <lease_expires_epoch>0</lease_expires_epoch>
        <lease_expires_utc/>
        <life_expectancy>0</life_expectancy>
        <purchasing_contact/>
        <applecare_id/>
        <attachments/>
    </purchasing>
</mobile_device>
    """

def create_site_config(site_name):
    return f"""
<site>
    <name>{site_name}</name>
</site>
    """


def parse_id_from_response(resp_text) -> str:
    start = "<id>"
    end = "</id>"
    return parse_tag_contents(start, end, resp_text)


def parse_tag_contents(start_tag, end_tag, resp_text):
    start_tag_length = len(start_tag)
    start_tag_index = resp_text.index(start_tag)
    end_tag_index = resp_text.index(end_tag)
    offset_start_tag_index = start_tag_index + start_tag_length

    return resp_text[offset_start_tag_index : end_tag_index]


def create_site(site_name):
    site_config = create_site_config(site_name)
    site_id = send_create(instance.classic.sites, site_config, "sites")
    return site_id


def create_mobile_devices(site_id, amount):
    device_ids = []
    for i in range (0, amount):
        device_name = f"tf-testing-{TESTING_ID}-{RANDOM_NUMBER}-{i}"
        device_config = create_mobile_device_config(device_name, site_id, SITE_NAME)
        device_id = send_create(instance.classic.mobile_devices, device_config, "mobile_devices")
        device_ids.append(device_id)
    return device_ids


def send_create(instance_object, payload, type_string):
    resp = instance_object.create(payload)
    if resp.ok:
        resp_text = resp.text
        object_id = parse_id_from_response(resp_text)
        logger.info(f"Successfully CREATED {type_string} id:{object_id}")
    else:
        logger.warning(f"FAILED to CREATE {type_string}")
    return object_id


def write_ids_to_data_source(site_id, device_ids):
    data_object = {
        "mobile_devices": device_ids,
            "site": site_id
        }
    full_path = "testing/data_sources/site_and_mobile_device_ids.json"
    file = Path(full_path)
    # If the file path doesnt exist, the next line facilitates its creation
    file.parent.mkdir(parents=True, exist_ok=True)
    data_json = json.dumps(data_object)
    file.write_text(data_json)


site_id = create_site(SITE_NAME)
device_ids = create_mobile_devices(site_id, DEVICE_COUNT)
write_ids_to_data_source(site_id=site_id, device_ids=device_ids)