package main

import "github.com/urfave/cli"

var AppConfigFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "node-name",
		Value:       "node",
		Usage:       "The name of the node that will be broadcasted among the network",
		EnvVar:      "NODE_NAME",
		Destination: &conf.NodeName,
	},
	cli.StringFlag{
		Name:        "listenaddr",
		Value:       "127.0.0.1:10000",
		Usage:       "The address on which the node will begin listening to for requests",
		EnvVar:      "NODE_LISTEN_ADDR",
		Destination: &conf.NodeAddr,
	},
	cli.StringFlag{
		Name:        "service-discover-addr",
		Value:       "127.0.0.1:8500",
		Usage:       "The address used for peer service discovery",
		EnvVar:      "SERVICE_DISCOVERY_ADDR",
		Destination: &conf.ServiceDiscoveryAddress,
	},
}
