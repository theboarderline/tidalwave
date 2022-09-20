package tidalwave

type ClusterCreater interface {
	Create() error
}

func CreateCluster(c ClusterCreater) error {
	err := c.Create()
	if err != nil {
		return err
	}
	return nil
}
