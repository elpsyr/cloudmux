package cloudprovider

import "yunion.io/x/onecloud/pkg/apis"

type CfelSManagedVMCreateConfig struct {
	SManagedVMCreateConfig
	IsolatedDevice       []*IsolatedDeviceConfig
	Networks             []Network
	BaremetalDiskConfigs []*BaremetalDiskConfig
	EipBw                int
	EipAutoDellocate     bool
	Count                int // 数量
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
type BaremetalDiskConfig struct {
	//Index int `json:"index"`
	// disk type
	Type string `json:"type"` // ssd / rotate
	// raid config
	Conf         string  `json:"conf"`  // raid配置
	Count        int64   `json:"count"` // 连续几块
	Range        []int64 `json:"range"` // 指定几块
	Splits       string  `json:"splits"`
	Size         []int64 `json:"size"` //
	Adapter      *int    `json:"adapter,omitempty"`
	Driver       string  `json:"driver"`
	Cachedbadbbu *bool   `json:"cachedbadbbu,omitempty"`
	Strip        *int64  `json:"strip,omitempty"`
	RA           *bool   `json:"ra,omitempty"`
	WT           *bool   `json:"wt,omitempty"`
	Direct       *bool   `json:"direct,omitempty"`
}

type IsolatedDeviceInfo struct {
	DevType        string `json:"dev_type,omitempty"`
	Model          string `json:"model,omitempty"`
	VendorDeviceId string `json:"vendor,omitempty"`
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
	GroupBy []GroupBy
}

type GroupBy struct {
	Type string
	Params []string
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

type GetNetworkOptions struct {
	ZoneId       string `json:"zone_id,omitempty"`
	WithUserMeta bool   `json:"with_user_meta,omitempty"`
	ServerType   string `json:"server_type,omitempty"`
	WireId       string `json:"wire_id,omitempty"`
	VpcId        string `json:"vpc_id,omitempty"`
}

type CfelSManagedVMRebuildRootConfig struct {
	SManagedVMRebuildRootConfig
	ResetPassword bool `json:"reset_password"`
	// 重置指定密码
	Password string `json:"password"`

	
	AutoStart bool `json:"auto_start"`
	DeployTelegraf bool `json:"deploy_telegraf"`
}
// ServerSSHInfo
// copy from ServerRemoteConsoleResponse
type ServerSSHInfo struct {
	AccessUrl     string `json:"access_url"`
	ConnectParams string `json:"connect_params"`
	Session       string `json:"session,omitempty"`

	apis.Meta
}
