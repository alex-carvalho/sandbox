# Ansible Collections

Ansible Collections are a distribution format for bundling and sharing playbooks, roles, modules, and plugins. They allow teams to package reusable automation content under a `namespace.collection` naming convention and install third-party collections from Ansible Galaxy.

## Project Structure

```
.
├── playbooks/
│   └── site.yml                          # Uses all three collections (FQCN syntax)
├── collections/ansible_collections/
│   └── myorg/utils/                      # Custom local collection
│       ├── galaxy.yml                    # Collection metadata
│       ├── meta/runtime.yml
│       └── plugins/modules/
│           └── system_info.py            # Custom module
├── requirements.yml                      # Galaxy collections to install
├── inventory
└── docker-compose.yml
```

## Collections Used

| Collection | Source | Modules used |
|---|---|---|
| `community.crypto` | Ansible Galaxy | `openssl_privatekey`, `openssl_csr`, `x509_certificate`, `x509_certificate_info` |

## Usage

```bash
make up      # Build images and start containers
make init    # Install Galaxy collections from requirements.yml
make run     # Run the playbook
make list    # List all installed collections
```

## Testing

Access the deployed pages over HTTPS (self-signed cert — accept the browser warning):
- Node 1: https://localhost:8443
- Node 2: https://localhost:8444

Or via curl skipping cert verification:
```bash
curl -k https://localhost:8443
curl -k https://localhost:8444
```
