# Linux Services with systemctl

## Overview

`systemctl` is the primary command-line tool for managing systemd services on modern Linux systems. Systemd is the system and service manager that replaces older init systems like SysVinit.

## Basic Service Management Commands

### Viewing Service Status
```bash
# Check the status of a specific service
systemctl status nginx

# List all active services
systemctl list-units --type=service --state=running

# List all services (active and inactive)
systemctl list-units --type=service --all

# Check if a service is enabled at boot
systemctl is-enabled nginx
```

### Starting and Stopping Services
```bash
# Start a service
sudo systemctl start nginx

# Stop a service
sudo systemctl stop nginx

# Restart a service
sudo systemctl restart nginx

# Reload service configuration (without restarting)
sudo systemctl reload nginx

# Reload or restart if necessary
sudo systemctl reload-or-restart nginx
```

### Enabling Services at Boot
```bash
# Enable a service to start at boot
sudo systemctl enable nginx

# Disable a service from starting at boot
sudo systemctl disable nginx

# Enable and start a service in one command
sudo systemctl enable --now nginx
```

## Service States

| State | Description |
|-------|-------------|
| `active (running)` | Service is currently running |
| `active (exited)` | Service ran successfully and exited |
| `inactive (dead)` | Service is stopped |
| `failed` | Service failed to start |
| `enabled` | Service starts at boot |
| `disabled` | Service does not start at boot |

## Inspecting Service Configuration

```bash
# View the service unit file
systemctl cat nginx

# Show the full configuration
systemctl show nginx

# Display specific properties
systemctl show -p MainPID nginx
systemctl show -p ExecStart nginx

# Check service dependencies
systemctl list-dependencies nginx
```

## Logs and Debugging

```shell
# View service logs
journalctl -u nginx

# View recent logs
journalctl -u nginx -n 50

# Follow logs in real-time
journalctl -u nginx -f

# View logs since boot
journalctl -u nginx -b

# View logs for a specific time period
journalctl -u nginx --since "2 hours ago"
```


Example service configuration:
```shell
cat /usr/lib/systemd/system/app.service
```

```
[Unit]
Description=My python web application

[Service]
ExecStart=/usr/bin/python3 /opt/code/my_app.py
ExecStartPre=/bin/bash /opt/code/configure_db.sh
ExecStartPost=/bin/bash /opt/code/email_status.sh
Restart=always

[Install]
WantedBy=multi-user.target
```