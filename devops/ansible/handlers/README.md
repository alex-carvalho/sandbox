# Ansible handlers 

"Ansible handlers are used to trigger actions, typically at the end of a play, only if a preceding task has made changes to the system. This mechanism ensures that actions dependent on configuration changes, such as restarting a service, are performed efficiently and only when necessary, maintaining the idempotency of the playbooks."


handlers can be trigger using `notify`and also forced to be executed using `meta: flush_handlers`.

## Run

```bash
make run
make check
```
