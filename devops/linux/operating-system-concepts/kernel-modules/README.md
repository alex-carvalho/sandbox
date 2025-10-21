Create a kernel module and load it to kernel

```shell
docker build -t kernel-modules .

# --privileged: Grants the container all capabilities, including CAP_SYS_MODULE which is required for loading modules.
# --cap-add=ALL: Explicitly adds all capabilities, ensuring CAP_SYS_MODULE is present.
docker run --rm -it  --privileged --cap-add=ALL kernel-modules /bin/bash

# insider container
make

# load the module
insmod simple.ko

# check module loaded 
lsmod | grep simple

# remove module
rmmod simple

# check logs
dmesg | grep "SIMPLE:"
```
