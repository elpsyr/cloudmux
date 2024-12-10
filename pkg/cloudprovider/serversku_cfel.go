package cloudprovider

type CfelSkuPriceOptions struct {
	ZoneId       string
	InstanceType string
	DiskType     string
	DiskSize     int
	ImageId      string
	FeeUnit      string //  Year：年。 Month：月 。
	ChargeType   string // 付费类型
	Duration     int    // 订购时长
	Quantity     int    // 数量
}
