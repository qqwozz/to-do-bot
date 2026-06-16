# Multi-stage build
FROM ubuntu:22.04 AS builder

RUN apt-get update && apt-get install -y \
    build-essential cmake pkg-config \
    libsqlite3-dev libssl-dev curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build
COPY CMakeLists.txt include/ src/ main.cpp ./
RUN mkdir -p build && cd build && cmake .. && make -j$(nproc)

# Runtime
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y \
    libsqlite3-0 curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /build/build/todo-backend .
RUN mkdir -p /app/data
EXPOSE 8081
CMD ["./todo-backend"]
