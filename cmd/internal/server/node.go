package server

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	hs "github.com/cpurta/p2p-grpc/proto-go"
	"github.com/hashicorp/consul/api"
)

const KeyPrefix = "grpc-p2p-"

type Node struct {
	Name string
	Addr string

	SDAddress string
	SDKV      api.KV

	Peers map[string]hs.HelloServiceClient
}

func (node *Node) SayHello(ctx context.Context, stream *hs.HelloRequest) (*hs.HelloReply, error) {
	return &hs.HelloReply{Message: "Hello from " + node.Name}, nil
}

func (node *Node) StartListening() {
	lis, err := net.Listen("tcp", node.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer() // n is for serving purpose

	hs.RegisterHelloServiceServer(grpcServer, node)
	reflection.Register(grpcServer)

	// start listening
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (node *Node) registerService() error {
	config := api.DefaultConfig()
	config.Address = node.SDAddress
	consul, err := api.NewClient(config)
	if err != nil {
		log.Println("Unable to contact Service Discovery:", err)
		return err
	}

	kv := consul.KV()
	pair := &api.KVPair{Key: KeyPrefix + node.Name, Value: []byte(node.Addr)}
	_, err = kv.Put(pair, nil)
	if err != nil {
		log.Println("Unable to register with Service Discovery:", err)
		return err
	}

	// store the kv for future use
	node.SDKV = *kv

	log.Println("Successfully registered with Consul.")
	return nil
}

func (node *Node) Start() error {
	node.Peers = make(map[string]hs.HelloServiceClient)

	go node.StartListening()

	if err := node.registerService(); err != nil {
		return err
	}

	for {
		// TODO: pull messages from centralized channel
		// and broadcast those to all known peers
		// let just sleep for now
		node.BroadcastMessage("Hello from " + node.Name)
		time.Sleep(time.Minute * 1)
	}
}

func (node *Node) BroadcastMessage(message string) {
	// get all nodes -- inefficient, but this is just an example
	kvpairs, _, err := node.SDKV.List(KeyPrefix, nil)
	if err != nil {
		log.Println("Error getting keypairs from service discovery:", err)
		return
	}

	for _, kventry := range kvpairs {
		if strings.Compare(kventry.Key, KeyPrefix+node.Name) == 0 {
			// ourself
			continue
		}
		if node.Peers[kventry.Key] == nil {
			fmt.Println("New member: ", kventry.Key)
			// connection not established previously
			node.SetupClient(kventry.Key, string(kventry.Value))
		}
	}
}

func (node *Node) SetupClient(name string, addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Unable to connect to %s: %v", addr, err)
	}

	defer conn.Close()

	node.Peers[name] = hs.NewHelloServiceClient(conn)

	response, err := node.Peers[name].SayHello(context.Background(), &hs.HelloRequest{Name: node.Name})
	if err != nil {
		log.Printf("Error making request to %s: %v", name, err)
	}

	log.Printf("Greeting from other node: %s", response.GetMessage())
}
