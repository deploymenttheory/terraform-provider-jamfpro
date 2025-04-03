import os
import jamfpy
import random
import uuid
import json
from pathlib import Path
from dotenv import load_dotenv
load_dotenv()

logger = jamfpy.get_logger(name="site_computer_setup", level=20)
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

site_name = "tf-testing-site"

def create_computer_config(computer_name,site_id, site_name):
    return f"""
<computer>
    <general>
        <name>{computer_name}</name>
        <serial_number>{uuid.uuid4()}</serial_number>
        <udid>{uuid.uuid4()}</udid>
        <barcode_1/>
        <barcode_2/>
        <asset_tag/>
        <remote_management>
            <managed>true</managed>
            <management_username>jamfadmin</management_username>
            <management_password>string</management_password>
        </remote_management>
        <site>
            <id>-{site_id}</id>
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
        <applecare_id>test</applecare_id>
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
        <os_applecare_id/>
        <os_maintenance_expires/>
        <attachments/>
    </purchasing>
    <extension_attributes>
        <extension_attribute>
            <id>2</id>
            <value/>
        </extension_attribute>
    </extension_attributes>
</computer>
    """
    return

def create_site_config(site_name):
    return f"""
<site>
    <name>{site_name}</name>
</site>
    """

def parse_id_from_response(resp_text) -> str: 
    start = "<id>"
    end = "</id>"
    return resp_text[resp_text.index(start) + len(start): resp_text.index(end)]

def create_site(site_name):
    site_config = create_site_config(site_name)
    resp = instance.classic.sites.create(site_config)
    resp.raise_for_status()
    resp_text = resp.text
    site_id = parse_id_from_response(resp_text)
    logger.info(f"Sucessfully created site {site_name} id:{site_id}")
    return site_id

def create_computers(site_id):
    computer_ids = []
    random_number = random.randint(0,9999)
    computer_name_1 = f"tf-testing-{random_number}-1"
    computer_name_2 = f"tf-testing-{random_number}-2"
    computer_config_1 = create_computer_config(computer_name_1, site_id, site_name)
    computer_config_2 = create_computer_config(computer_name_2, site_id, site_name)
    computer_configs = [computer_config_1, computer_config_2]
    for config in computer_configs:
        resp = instance.classic.computers.create(config)
        resp.raise_for_status()
        response_text = resp.text
        computer_id = parse_id_from_response(response_text)
        logger.info(f"Sucessfully created computer id:{computer_id}")
        computer_ids.append(int(computer_id))
    return computer_ids



def write_ids_to_data_source(site_id, computer_ids):
    data_object = {
        "computers":computer_ids,
            "site":site_id
            }
    full_path = "../data_sources/site_and_computer_ids.json"
    file = Path(full_path)
    file.parent.mkdir(parents=True, exist_ok=True)
    data_json = json.dumps(data_object)
    file.write_text(data_json)

site_id = create_site(site_name)
computer_ids = create_computers(site_id)
write_ids_to_data_source(site_id=site_id, computer_ids=computer_ids)