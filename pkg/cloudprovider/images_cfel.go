package cloudprovider

type CfelSImageCreateOption struct {
	IsProtected    bool   `json:"protected,omitempty"`
	ImageId        string `json:"image_id,omitempty"`
	ExternalId     string `json:"external_id,omitempty"`
	ImageName      string `json:"image_name,omitempty"`
	CopyFrom       string `json:"copy_from,omitempty"`
	Description    string `json:"description,omitempty"`
	MinDiskMb      int    `json:"min_disk_mb,omitempty"`
	MinRamMb       int    `json:"min_ram_mb,omitempty"`
	Checksum       string `json:"checksum,omitempty"`
	OsType         string `json:"os_type,omitempty"`
	OsArch         string `json:"os_arch,omitempty"`
	OsDistribution string `json:"os_distribution,omitempty"`
	OsVersion      string `json:"os_version,omitempty"`
	OsFullVersion  string `json:"os_full_version,omitempty"`
}

type CfelResetGuestPasswordOption struct {
	GuestID       string `json:"guest_id,omitempty"`
	ResetPassword bool   `json:"reset_password,omitempty"`
	AutoStart     bool   `json:"auto_start,omitempty"`
	Password      string `json:"password,omitempty"`
	UserName      string `json:"username,omitempty"`
}

type CfelSetImageUserTag struct {
	ImageId string
	Tags map[string]string
}

type CfelChangeSettingOption struct {
	Desc     string `json:"description,omitempty"`
	Name     string `json:"name,omitempty"`
	HostName string `json:"hostname,omitempty"`
}