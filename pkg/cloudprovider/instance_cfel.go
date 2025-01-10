package cloudprovider

import (
	"yunion.io/x/onecloud/pkg/apis"
)

const (
	InstanceChargeTypeTag      = "instanceChargeType"
	InstanceChargeTypeSpotPaid = "SpotPaid"
	InstanceChargeTypePostPaid = "PostPaid"
	InstanceChargeTypePrePaid  = "PrePaid"
)

type CfelServerRebuildRootInput struct {
	apis.Meta

	// swagger: ignore
	Image string `json:"image" yunion-deprecated-by:"image_id"`
	// 关机且停机不收费情况下不允许重装系统
	// 镜像 id
	// required: true
	ImageId string `json:"image_id"`
	// swagger: ignore
	// Keypair string `json:"keypair" yunion-deprecated-by:"keypair_id"`
	// 秘钥Id
	// KeypairId     string `json:"keypair_id"`
	// ResetPassword *bool  `json:"reset_password"`
	// Password      string `json:"password"`

	AutoStart *bool `json:"auto_start"`

	AllDisks *bool `json:"all_disks"`

	CfelServerDeployInputBase
}

type CfelServerDeployInputBase struct {
	// swagger: ignore
	Keypair string `json:"keypair,omitempty" yunion-deprecated-by:"keypair_id"`
	// 秘钥Id
	KeypairId string `json:"keypair_id,omitempty"`

	// 清理指定公钥
	// 若指定的秘钥Id和虚拟机的秘钥Id不相同, 则清理旧的公钥
	DeletePublicKey string `json:"delete_public_key,omitempty"`
	// 解绑当前虚拟机秘钥, 并清理公钥信息
	DeleteKeypair bool `json:"__delete_keypair__,omitempty"`
	// 生成随机密码, 优先级低于password
	ResetPassword bool `json:"reset_password,omitempty"`
	// 重置指定密码
	Password string `json:"password,omitempty"`
	// swagger: ignore
	LoginAccount string `json:"login_account,omitempty"`

	// swagger: ignore
	Restart bool `json:"restart,omitempty"`

	// swagger: ignore
	//DeployConfigs []*DeployConfig `json:"deploy_configs"`
	// swagger: ignore
	DeployTelegraf bool `json:"deploy_telegraf,omitempty"`

	UserData string `json:"user_data,omitempty"`
}

type CfelSManagedVMCreateConfig struct {
	SManagedVMCreateConfig
	IsolatedDevice       []*IsolatedDeviceConfig
	Networks             []Network
	BaremetalDiskConfigs []*BaremetalDiskConfig
	EipBw                int
	EipAutoDellocate     bool
	Count                int    // 数量
	PreferHostId         string // 调度使用指定宿主机
	PreferZoneId         string // 调度使用指定宿主机
	Machine              string // emulate: pc, q35
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
	GroupBy  []GroupBy
}

type GroupBy struct {
	Type   string
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
	ZoneId       string   `json:"zone_id,omitempty"`
	WithUserMeta bool     `json:"with_user_meta,omitempty"`
	ServerType   string   `json:"server_type,omitempty"`
	WireId       string   `json:"wire_id,omitempty"`
	VpcId        string   `json:"vpc_id,omitempty"`
	Ids          []string `json:"ids,omitempty"`
}

type CfelSManagedVMRebuildRootConfig struct {
	SManagedVMRebuildRootConfig
	ResetPassword bool `json:"reset_password"`
	// 重置指定密码
	Password string `json:"password"`

	AutoStart      bool `json:"auto_start"`
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

type CfelResetGuestPasswordOption struct {
	GuestID       string `json:"guest_id,omitempty"`
	ResetPassword bool   `json:"reset_password,omitempty"`
	AutoStart     bool   `json:"auto_start,omitempty"`
	Password      string `json:"password,omitempty"`
	UserName      string `json:"username,omitempty"`
}

type CfelExecCmdOption struct {
	Ip       string
	Port     int
	Cmd      string
	User     string
	Password string
}
