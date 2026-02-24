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

| Collection | Source | Module used |
|---|---|---|
| `community.general` | Ansible Galaxy | `ini_file` — writes INI config files |
| `ansible.posix` | Ansible Galaxy | `acl` — manages file ACL permissions |
| `myorg.utils` | Local | `system_info` — custom module, returns hostname/platform/Python version |

## Usage

```bash
make up      # Build images and start containers
make init    # Install Galaxy collections from requirements.yml
make run     # Run the playbook
make list    # List all installed collections
```

## Testing

Access the deployed pages:
- Node 1: http://localhost:8001
- Node 2: http://localhost:8002
