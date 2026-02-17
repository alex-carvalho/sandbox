FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    build-essential \
    libgflags-dev \
    libsnappy-dev \
    zlib1g-dev \
    libbz2-dev \
    liblz4-dev \
    libzstd-dev \
    git \
    cmake \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /opt
RUN git clone --depth 1 --branch v10.10.1 https://github.com/facebook/rocksdb.git && \
    cd rocksdb && \
    make static_lib shared_lib -j$(nproc) && \
    make install && \
    ldconfig

WORKDIR /app
COPY client.cpp .

RUN g++ -std=c++20 client.cpp -o client -lrocksdb -lpthread -lz -lbz2 -llz4 -lzstd -lsnappy

CMD ["tail", "-f", "/dev/null"]
