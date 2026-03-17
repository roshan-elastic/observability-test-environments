"""Script to configure fleet."""

import os

import requests
from requests.auth import HTTPBasicAuth

es_url = os.getenv("ELASTICSEARCH_HOST")
kb_url = os.getenv("KIBANA_HOST")
es_user = os.getenv("ELASTICSEARCH_USERNAME")
es_pass = os.getenv("ELASTICSEARCH_PASSWORD")
basic_auth = HTTPBasicAuth(es_user, es_pass)
default_headers = {
    "Content-Type": "application/json",
    "kbn-xsrf": "true"
}

wait = True
while wait:
    try:
        res = requests.get(
            url=es_url,
            auth=basic_auth,
            headers=default_headers,
            timeout=10
        )
        if res.status_code == 200:
            wait = False
    except requests.exceptions.RequestException:
        print("Waiting for Elasticsearch")

token_url = f"{es_url}/_security/oauth2/token"
data_json = {"grant_type": "client_credentials"}
res = requests.post(
    url=token_url,
    auth=basic_auth,
    json=data_json,
    headers=default_headers,
    timeout=10
)
es_access_token = res.json()['access_token']

service_toke_url = f"{es_url}/_security/service/elastic/fleet-server/credential/token"
data_json = {"grant_type": "client_credentials"}
default_headers['Authorization'] = f"Bearer {es_access_token}"
res = requests.post(
    url=service_toke_url,
    headers=default_headers,
    timeout=10
)
fleet_service_token = res.json()['token']['value']

print(f"ES_ACCESS_TOKEN={es_access_token}")
print(f"FLEET_SERVICE_TOKEN={fleet_service_token}")

with open("/out/environment", "w", encoding="UTF-8") as f:
    f.write(f"export ES_ACCESS_TOKEN={es_access_token}\n")
    f.write(f"export FLEET_SERVER_SERVICE_TOKEN={fleet_service_token}\n")
    f.write(f"export KIBANA_FLEET_SERVICE_TOKEN={fleet_service_token}\n")
