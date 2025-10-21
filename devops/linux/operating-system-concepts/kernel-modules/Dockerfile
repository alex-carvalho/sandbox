FROM ubuntu:latest

RUN apt update && apt install -y build-essential kmod linux-headers-$(uname -r)

COPY src/ /src/

WORKDIR /src

CMD ["bash"]
