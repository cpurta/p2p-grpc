FROM golang:latest

RUN apt-get update && \
    mkdir -p /app

ADD ./bin/p2p-grpc /app/p2p-grpc

RUN chmod +x /app/p2p-grpc

EXPOSE 10000
ENTRYPOINT ["/app/p2p-grpc"]
