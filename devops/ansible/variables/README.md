# Ansible variables

Ansible testing environment with one master and two worker nodes running Apache.

## Testing

```bash
docker-compose up 

# Wait SSH initialization, then run the playbook
docker exec -it ansible-master ansible-playbook /etc/ansible/playbooks/apache.yml

# Test web servers
# Node 1 (Red): http://localhost:8001
# Node 2 (Blue): http://localhost:8002
```

## Structure

- `playbooks/apache.yml` - Apache installation playbook
- `host_vars/node1` - Node 1 variables (red background)
- `host_vars/node2` - Node 2 variables (blue background)
- `assets/index.html` - Web page template
- `inventory` - Ansible inventory


## Variables

Used variables To change the background colors: `host_vars/node1`  and `host_vars/node2` 
