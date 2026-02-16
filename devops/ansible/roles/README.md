# Ansible Roles

This project demonstrates Ansible automation using roles for modular and reusable infrastructure management.

## Project Structure

```
.
├── playbooks/           # Ansible playbooks
│   └── apache.yml      # Apache web server playbook
├── roles/              # Ansible roles
│   └── apache/         # Apache role
│       ├── defaults/   # Default variables
│       ├── handlers/   # Handlers (e.g., restart service)
│       ├── tasks/      # Main tasks
│       └── templates/  # Jinja2 templates
├── host_vars/          # Host-specific variables
├── inventory           # Inventory file
└── docker-compose.yml  # Test environment
```

## Apache Role

The Apache role installs and configures Apache web server with the following features:
- Package cache update
- Network utilities installation (curl, wget)
- Apache2 installation
- Custom index.html deployment
- Service management with handlers

### Usage

Run the Apache playbook:
```bash
# Start the test environment
docker-compose up -d

# Execute the playbook
docker exec ansible-master ansible-playbook /etc/ansible/playbooks/apache.yml -i /etc/ansible/inventory
```

### Variables

Default variables (defined in `roles/apache/defaults/main.yml`):
- `apache_port`: 80
- `document_root`: /var/www/html
- `network_utilities`: [curl, wget]

Host-specific variables can be set in `host_vars/`.

### Testing

Access the web servers:
- Node 1: http://localhost:8001
- Node 2: http://localhost:8002

Each node displays a custom page with its unique configuration.