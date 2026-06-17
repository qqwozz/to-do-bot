# Multi-stage build
FROM ubuntu:22.04 AS builder

RUN apt-get update && apt-get install -y \
    build-essential cmake pkg-config \
    libsqlite3-dev libssl-dev libasio-dev curl git \
    && rm -rf /var/lib/apt/lists/*

RUN git clone --depth 1 https://github.com/CrowCpp/Crow.git /tmp/crow && \
    cp -r /tmp/crow/include/crow /usr/local/include/crow && \
    cp /tmp/crow/include/crow.h /usr/local/include/crow.h && \
    rm -rf /tmp/crow

WORKDIR /build
COPY CMakeLists.txt ./
COPY include/ include/
COPY src/ src/
COPY main.cpp ./
RUN mkdir -p build && cd build && \
    cmake .. -DCROW_INCLUDE_DIR=/usr/local/include && \
    make -j$(nproc)

# Runtime
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y \
    libsqlite3-0 libssl3 curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /build/build/todo-backend .
RUN mkdir -p /app/data
EXPOSE 8081
CMD ["./todo-backend"]
