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
