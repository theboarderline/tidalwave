package google

import (
	serviceusage "cloud.google.com/go/serviceusage/apiv1"
	"context"
	"fmt"
	"github.com/kyokomi/emoji/v2"

	//"github.com/kyokomi/emoji/v2"
	serviceusagepb "google.golang.org/genproto/googleapis/api/serviceusage/v1"
)

type Apis struct {
	Services []string
}

// RequiredApis is a static declaration of required Apis - should be vetted to be comprehensive
var RequiredApis = Apis{
	Services: []string{
		"cloudkms.googleapis.com",
		"compute.googleapis.com",
		"container.googleapis.com",
		"iam.googleapis.com",
	},
}

func (r *Controlplane) EnableApis() (error) {
	ctx := context.Background()
	c, err := serviceusage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	req  := &serviceusagepb.BatchEnableServicesRequest{
		Parent: fmt.Sprintf("projects/%s", r.ProjectNumber),
		ServiceIds: RequiredApis.Services,
	}

	operation, err := c.BatchEnableServices(ctx, req)
	if err != nil {
		return err
	}
	resp, err := operation.Wait(ctx)
	if err != nil {
		return err
	}
	for _, service := range resp.Services {
		emoji.Printf(":check_mark_button: Enabled %s\n", service.Name)
	}

	return nil
}
