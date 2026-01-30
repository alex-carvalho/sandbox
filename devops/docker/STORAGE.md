# Docker Storage & Bind Mounts

## Storage Types

```bash
# bind directory/file

# Read-write
docker run --mount type=bind,source=/host/data,target=/app/data nginx

# Read-only
docker run --mount type=bind,source=/host/data,target=/app/data,readonly=true nginx


### volumes

# Create volume
docker volume create mydata

# Use named volume
docker run --mount type=volume,source=mydata,target=/data nginx

## tmpfs

#In-memory temporary storage (lost on container stop).
docker run --mount type=tmpfs,target=/tmp nginx

# With size limit
docker run --mount type=tmpfs,target=/cache,tmpfs-size=256M nginx

# With permissions
docker run --mount type=tmpfs,target=/run,tmpfs-size=512M,tmpfs-mode=1777 nginx
```

## Comparison

| Type | Location | Persistent | Speed | Use Case |
|------|----------|-----------|-------|----------|
| Bind | Host filesystem | Yes | Medium | Development, configs |
| Volume | Docker area | Yes | Fast | Production, databases |
| Tmpfs | Memory | No | Fastest | Temp data, secrets |
