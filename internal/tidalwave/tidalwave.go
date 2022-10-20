package tidalwave

import (
	"strings"
)

// ClusterCreater provides cluster creation
type ClusterCreater interface {
	Create() error
	EnableApis() error
}

//CreateApis implements the ClusterCreater interface to enable apis
func CheckApis(c ClusterCreater) (error) {
	if err := c.EnableApis(); err != nil {
		return err
	}
	return nil
}

// CreateCluster creates a cluster and dependencies
func CreateCluster(c ClusterCreater) error {
	err := c.Create()
	if err != nil {
		return err
	}
	return nil
}

// ClusterDeleter provides cluster deletion
type ClusterDeleter interface {
	Delete() error
}

// DeleteCluster deletes a cluster and dependencies
func DeleteCluster(c ClusterDeleter) error {
	err := c.Delete()
	if err != nil {
		return err
	}
	return nil
}

// ClusterUpdater provides cluster updates
type ClusterUpdater interface {
	Update() error
}

// UpdateCluster updates a cluster and dependencies
func UpdateCluster(c ClusterUpdater) error {
	err := c.Update()
	if err != nil {
		return err
	}
	return nil
}

// NameFormatter removes -controlplane from cluster resource names if it exists
func NameFormatter(s string) string {
	return strings.TrimSuffix(s, "-controlplane")
}
