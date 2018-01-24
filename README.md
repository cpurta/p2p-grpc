# p2p-grpc

Experimental peer-to-peer network using Google's gRPC. This is a simple network whose
messages are just hello messages from the nodes to other nodes. It relies on [hashicorp
consul](https://github.com/hashicorp/consul) for service discovery of other nodes and
for node key/value storage.

## Requirements

You will need to either have consul installed locally or you can pull the docker image
so that you can have a consul agent that local nodes can connect to.

It is recommended that you also have [Go](https://golang.org/dl/) installed (1.9+)

## Quick Start

You will need to build the binary by running the following command:
```
$ go build -o ./bin/p2p-grpc ./cmd/p2p-grpc
```

In another terminal tab/window you will need to start consul:
```
$ consul agent -dev
```

You can then start the first node:
```
$ ./bin/p2p-grpc --node-name node-1 --listenaddr 127.0.0.1:10000 --service-discover-addr=127.0.0.1:8500
```

and lets start another:
```
$ ./bin/p2p-grpc --node-name node-2 --listenaddr 127.0.0.1:10001 --service-discover-addr=127.0.0.1:8500
```

## Issues

There are currently some small issues that have been encountered with using consul as
the key/value storage system for the network. These will be addressed later in the future
as I mess around with using consul or perhaps add key-value storage on a node basis.

## LICENSE

MIT
