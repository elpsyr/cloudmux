package cloudprovider

type CfelSManagedVMCreateConfig struct {
	SManagedVMCreateConfig
	IsolatedDevice []*IsolatedDeviceConfig
}

type IsolatedDeviceConfig struct {
	DevType      string `json:"dev_type"`
	DiskIndex    int    `json:"disk_index"`
	Id           string `json:"id"`
	Index        int    `json:"index"`
	Model        string `json:"model"`
	NetworkIndex int    `json:"network_index"`
	Vendor       string `json:"vendor"`
	WireId       string `json:"wire_id"`
}
