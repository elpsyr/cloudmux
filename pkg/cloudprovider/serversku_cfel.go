package cloudprovider

type CfelSkuPriceOptions struct {
	ZoneId       string
	InstanceType string
	SysDiskType  string
	SysDiskSize  int
	ImageId      string
	FeeUnit      string //  Year：年。 Month：月 。
	ChargeType   string
	Duration     int // 订购时长
	Quantity     int // 数量
}
