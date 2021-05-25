package main

import (
	gcp "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	mcp "istio.io/api/mcp/v1alpha1"
	networking "istio.io/api/networking/v1alpha3"

	"log"
	"net"
	"time"
)

var version = time.Now().Format("20060102150405")

func main() {
	log.SetFlags(log.Lshortfile)
	listener, err := net.Listen("tcp", ":9901")
	if err != nil {
		log.Println(err)
	}
	s := grpc.NewServer()

	t := Testmcp{}
	gcp.RegisterAggregatedDiscoveryServiceServer(s, t)
	reflection.Register(s)
	log.Println("begin serve")
	s.Serve(listener)
}

type Testmcp struct {
}

func (t Testmcp) StreamAggregatedResources(resourcesServer gcp.AggregatedDiscoveryService_StreamAggregatedResourcesServer) error {
	data := &networking.ServiceEntry{
		Hosts: []string{"a.b.c"},
		Ports: []*networking.Port{
			{
				Number:   20882,
				Name:     "say hi",
				Protocol: "TCP",
			},
		},
		Location:   networking.ServiceEntry_MESH_INTERNAL,
		Resolution: networking.ServiceEntry_STATIC,
		Endpoints: []*networking.WorkloadEntry{
			{
				Address: "1.1.1.1",
			},
			{
				Address: "2.2.2.2",
			},
			{
				Address: "3.3.3.3",
			},
		},
	}

	d, _ := types.MarshalAny(data)
	ptime, _ := types.TimestampProto(time.Now())
	res := &mcp.Resource{
		Metadata: &mcp.Metadata{
			Name:       "default/test",
			Version:    version,
			CreateTime: ptime,
		},
		Body: d,
	}

	b := proto.NewBuffer(nil)
	b.SetDeterministic(true)
	err := b.Marshal(res)
	if err != nil {
		log.Println(err)
	}
	resourcesServer.Send(&gcp.DiscoveryResponse{
		TypeUrl: "networking.istio.io/v1alpha3/ServiceEntry",
		Resources: []*any.Any{
			{
				TypeUrl: "type.googleapis.com/istio.mcp.v1alpha1.Resource",
				Value:   b.Bytes(),
			},
		},
	})
	time.Sleep(1 * time.Second * 3600)
	return nil
}

func (t Testmcp) DeltaAggregatedResources(resourcesServer gcp.AggregatedDiscoveryService_DeltaAggregatedResourcesServer) error {
	panic("implement me")
}
