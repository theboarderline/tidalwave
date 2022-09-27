package google

import (
	"context"
	"fmt"

	project "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "google.golang.org/genproto/googleapis/cloud/resourcemanager/v3"
)

func getProject(ctx context.Context, id string) (*resourcemanagerpb.Project, error) {
	client, err := project.NewProjectsClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	req := &resourcemanagerpb.SearchProjectsRequest{
		Query: fmt.Sprintf("name:%s", id),
	}
	resp := client.SearchProjects(ctx, req)
	return resp.Next()
}
