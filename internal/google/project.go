package google

import (
	"context"
	"fmt"
	"strings"

	resource "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "google.golang.org/genproto/googleapis/cloud/resourcemanager/v3"
)

func GetProjectNumber(id string) (*string, error) {
	ctx := context.Background()
	client, err := resource.NewProjectsClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	req := &resourcemanagerpb.GetProjectRequest{
		Name: fmt.Sprintf("projects/%s", id),
	}
	resp, err := client.GetProject(ctx, req)
	if err != nil {
		return nil, err
	}
	p := strings.Trim(resp.GetName(), "projects/")
	ptr := &p
	return ptr, nil
}
