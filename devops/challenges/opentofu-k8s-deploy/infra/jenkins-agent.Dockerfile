# Jenkins agent with Docker, Helm, Gradle, and Git
FROM jenkins/inbound-agent:latest-jdk21

USER root

# Install Docker Engine (full)
RUN apt-get update && \
    apt-get install -y apt-transport-https ca-certificates curl gnupg2 lsb-release software-properties-common && \
    curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add - && \
    echo "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list && \
    apt-get update && \
    apt-get install -y docker-ce

# Install Helm
RUN curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Install Gradle
RUN apt-get install -y wget unzip && \
    wget https://services.gradle.org/distributions/gradle-8.7-bin.zip -O /tmp/gradle.zip && \
    unzip /tmp/gradle.zip -d /opt && \
    ln -s /opt/gradle-8.7/bin/gradle /usr/bin/gradle

# Install Git
RUN apt-get install -y git

