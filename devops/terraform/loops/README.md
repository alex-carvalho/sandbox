# Terraform Loops & Docker

## Features Demonstrated

- `for_each` with map variables - service containers
- `for_each` with list transformation - enabled services filtering
- `count` for multiple instances - worker containers
- nested iteration using `for` + `flatten` - environment x service matrix
- `for` expressions for transformations - data normalization
- filtering with `if` conditions - enable/disable logic

## Real Docker Resources

- **Docker Network**: Bridge network for container communication
- **Service Containers**: Created via `for_each` with dynamic port mapping
- **Worker Containers**: Created via `count` with auto-generated names
- **Environment Variables**: Injected from Terraform variables
- **Labels**: Applied to containers for organization

## Requirements

Make sure Docker is running:

```bash
docker ps
```

## Run

```bash
terraform init
terraform plan
terraform apply -auto-approve
terraform output -json
docker ps
```

## Clean Up

```bash
terraform destroy -auto-approve
```

## Files

- `variables.tf`: input structures defining environments, services, and workers
- `main.tf`: Docker network, containers, and loop patterns
- `outputs.tf`: container details and connection strings
- `terraform.tfvars`: sample values for execution

## Example Commands

```bash
docker ps --filter "label=managed_by=terraform"
docker network ls
docker logs service-api
```
