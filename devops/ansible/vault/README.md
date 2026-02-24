# Ansible Vault

Ansible Vault encrypts sensitive data (passwords, API keys, secrets) directly inside YAML files so they can be safely committed to version control. 
Encrypted files are decrypted transparently at playbook runtime using a password file or prompt.


## Usage

```bash
make up      # Build images and start containers
make run     # Run the playbook (vault is auto-decrypted at runtime)
make view    # Print decrypted vault contents
make edit    # Edit the vault file interactively
```

> `group_vars/all/vault.yml` is committed to git already encrypted. The vault password lives in `.vault_pass` which is git-ignored.

## Testing

Access the deployed pages:
- Node 1: http://localhost:8001
- Node 2: http://localhost:8002
