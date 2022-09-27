package google

// import (
// 	"context"
// 	"net"
// 	"testing"

// 	"github.com/googleapis/gax-go/v2"
// 	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
// 	"google.golang.org/grpc"
// )

// type fakeNetworkServer struct {
// 	computepb.InsertNetworkRequest
// }

// func (f *fakeNetworkServer) Insert(ctx context.Context, req *computepb.InsertNetworkRequest, opts ...gax.CallOption) (*computepb.Network, error) {
// 	resp := &computepb.Network{
// 		Name: strPtr("myvpc"),
// 	}
// 	return resp, nil
// }

// func TestNetworkCreation(t *testing.T) {
// 	ctx := context.Background()

// 	fakeNetworkServer := &fakeNetworkServer{}

// 	l, err := net.Listen("tcp", "localhost:0")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	gsrv := grpc.NewServer()

// }
