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

resp = instance.classic.scripts.get_all()
resp.raise_for_status()
all_scripts = resp.json()["scripts"]
for i in all_scripts:
    del_resp = instance.classic.scripts.delete_by_id(i["id"])
    print(del_resp.text)
