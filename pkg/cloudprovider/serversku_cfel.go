package cloudprovider

type CfelSkuPriceOptions struct {
	InstanceType string
	SysDiskType  string
	SysDiskSize  int
	ImageId      string
	FeeUnit      string
	ChargeType   string
	Duration     int // 订购时长
	Quantity     int
}
