# Ansible Blue-Green Deployment

Ansible update proxy config to switch from one to another


## Usage

```bash
make setup
make deploy
make check
```

## Testing

```bash
curl http://localhost:8080                    # see current active version
curl -sI http://localhost:8080 | grep X-Active-Slot
```
