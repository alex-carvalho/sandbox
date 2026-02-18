#!/usr/bin/env python3
"""
Dynamic Inventory Script for Ansible POC.

Ansible calls this script with either:
  --list  -> returns all groups/hosts/vars as JSON
  --host <hostname> -> returns vars for a single host as JSON

In this POC the inventory is hard-coded to mirror the two Docker nodes
defined in docker-compose.yml. In a real scenario this data would come
from an external source (AWS API, CMDB, database, etc.).
"""

import json
import sys


NODES = {
    "node1": {"ansible_host": "node1", "ansible_port": 22, "http_port": 8001},
    "node2": {"ansible_host": "node2", "ansible_port": 22, "http_port": 8002},
}

INVENTORY = {
    "webservers": {
        "hosts": list(NODES.keys()),
        "vars": {
            "ansible_user": "root",
            "ansible_password": "testpass",
            "ansible_connection": "ssh",
        },
    },
    "_meta": {
        "hostvars": NODES,
    },
}


def main():
    if len(sys.argv) == 2 and sys.argv[1] == "--list":
        print(json.dumps(INVENTORY, indent=2))
    elif len(sys.argv) == 3 and sys.argv[1] == "--host":
        hostname = sys.argv[2]
        print(json.dumps(NODES.get(hostname, {}), indent=2))
    else:
        print(json.dumps({}))


if __name__ == "__main__":
    main()
