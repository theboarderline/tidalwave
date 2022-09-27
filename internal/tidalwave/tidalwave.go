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

type ClusterDeleter interface {
	Delete() error
}

func DeleteCluster(c ClusterDeleter) error {
	err := c.Delete()
	if err != nil {
		return err
	}
	return nil
}

type ClusterUpdater interface {
	Update() error
}

func UpdateCluster(c ClusterUpdater) error {
	err := c.Update()
	if err != nil {
		return err
	}
	return nil
}
