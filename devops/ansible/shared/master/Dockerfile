FROM ubuntu:26.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    ansible \
    openssh-client \
    sshpass \
    python3 \
    python3-pip \
    sudo \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/ansible/playbooks \
             /etc/ansible/host_vars \
             /etc/ansible/roles

WORKDIR /etc/ansible

CMD ["/bin/bash"]
