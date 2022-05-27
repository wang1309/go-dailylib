package flow

// GetPartitionSize returns the size in MB for each partition of
// the dataset. This is based on the hinted total size divided by
// the number of partitions.
func (d *Dataset) GetPartitionSize() int64 {
	return d.GetTotalSize() / int64(len(d.Shards))
}

// GetTotalSize returns the total size in MB for the dataset.
// This is based on the given hint.
func (d *Dataset) GetTotalSize() int64 {
	if d.Meta.TotalSize >= 0 {
		return d.Meta.TotalSize
	}
	var currentDatasetTotalSize int64
	for _, ds := range d.Step.InputDatasets {
		currentDatasetTotalSize += ds.GetTotalSize()
	}
	d.Meta.TotalSize = currentDatasetTotalSize
	return currentDatasetTotalSize
}
