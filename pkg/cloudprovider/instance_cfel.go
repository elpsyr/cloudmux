package cloudprovider

type CfelSManagedVMCreateConfig struct {
	SManagedVMCreateConfig
	IsolatedDevice []*IsolatedDeviceConfig
	Networks       []Network
}
type Network struct {
	NetworkId      string
	RequireTeaming bool
	Address        string
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

type IsolatedDeviceInfo struct {
	DevType string `json:"dev_type,omitempty"`
	Model   string `json:"model,omitempty"`
	VendorDeviceId  string `json:"vendor,omitempty"`
}

type MigrateForecastOption struct {
	GuestId         string `json:"guest_id,omitempty"`
	LiveMigrate     bool   `json:"live_migrate,omitempty"`
	SkipCpuCheck    bool   `json:"skip_cpu_check,omitempty"`
	SkipKernelCheck bool   `json:"skip_kernel_check,omitempty"`
	IsRescueMode    bool   `json:"is_rescue_mode,omitempty"`
}

type MonitorDataJSONOption struct {
	Measure  string
	Field    string
	GuestID  string
	Start    string
	End      string
	Interval string
}

type FilterOption struct {
	ShowFailReason string   `json:"show_fail_reason"`
	HostType       string   `json:"host_type"`
	Limit          int      `json:"limit"`
	HostStatus     string   `json:"host_status"`
	OsArch         string   `json:"os_arch"`
	Field          []string `json:"field"`
	FilterIds      string   `json:"filter"`
	Details        bool     `json:"details"`
}
