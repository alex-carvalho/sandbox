FROM ubuntu:26.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    ansible \
    openssh-client \
    openssh-server \
    ssh \
    sshpass \
    python3 \
    python3-pip \
    sudo \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /root/.ssh && \
    ssh-keygen -t rsa -N "" -f /root/.ssh/id_rsa && \
    cat /root/.ssh/id_rsa.pub >> /root/.ssh/authorized_keys && \
    chmod 600 /root/.ssh/authorized_keys

RUN mkdir -p /etc/ansible && \
    mkdir -p /etc/ansible/playbooks && \
    mkdir -p /etc/ansible/host_vars && \
    mkdir -p /etc/ansible/roles

WORKDIR /etc/ansible

CMD ["/bin/bash"]
