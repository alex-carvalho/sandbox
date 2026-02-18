FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    ansible \
    openssh-client \
    sshpass \
    python3 \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/ansible/playbooks

WORKDIR /etc/ansible

CMD ["/bin/bash"]
