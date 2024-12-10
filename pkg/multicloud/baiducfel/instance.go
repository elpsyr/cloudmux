package baiducfel

import (
	"context"
	"fmt"
	"strings"
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SInstance struct {
	multicloud.SInstanceBase
	BaiduTags
	multicloud.SBillingBase

	host   *SHost
	region *SRegion
	image  *SImage

	AutoRenew   bool   `json:"autoRenew"`
	CardCount   string `json:"cardCount"`
	CPUCount    int    `json:"cpuCount"`
	CreateTime  string `json:"createTime"`
	CreatedFrom string `json:"createdFrom"`
	// DedicatedHostID       string        `json:"dedicatedHostId"`
	DeletionProtection int `json:"deletionProtection"`
	// DeploysetList         []interface{} `json:"deploysetList"`
	Desc                  string       `json:"desc"`
	EniNum                string       `json:"eniNum"`
	ExpireTime            time.Time    `json:"expireTime"`
	Hostname              string       `json:"hostname"`
	ID                    string       `json:"id"`
	ImageID               string       `json:"imageId"`
	InstanceType          string       `json:"instanceType"`
	InternalIP            string       `json:"internalIp"`
	Ipv6                  string       `json:"ipv6"`
	IsomerismCard         string       `json:"isomerismCard"`
	LocalDiskSizeInGB     int          `json:"localDiskSizeInGB"`
	MemoryCapacityInGB    int          `json:"memoryCapacityInGB"`
	Name                  string       `json:"name"`
	NetworkCapacityInMbps int          `json:"networkCapacityInMbps"`
	NicInfo               SInstanceNic `json:"nicInfo"`
	NpuVideoMemory        string       `json:"npuVideoMemory"`
	PaymentTiming         string       `json:"paymentTiming"`
	PlacementPolicy       string       `json:"placementPolicy"`
	PublicIP              string       `json:"publicIp"`
	RoleName              string       `json:"roleName"`
	Spec                  string       `json:"spec"`
	Status                string       `json:"status"`
	SubnetID              string       `json:"subnetId"`
	Tags                  []BaiduTags  `json:"tags"`
	VpcID                 string       `json:"vpcId"`
	ZoneName              string       `json:"zoneName"`
	Volumes               []SDisk      `json:"volumes"`
}

const ServiceInstance = "instances"

// AllocatePublicIpAddress implements cloudprovider.ICloudVM.
// Subtle: this method shadows the method (SInstanceBase).AllocatePublicIpAddress of SInstance.SInstanceBase.
func (s *SInstance) AllocatePublicIpAddress() (string, error) {
	panic("unimplemented")
}

// AttachDisk implements cloudprovider.ICloudVM.
func (s *SInstance) AttachDisk(ctx context.Context, diskId string) error {
	var params = map[string]interface{}{
		"instanceId": s.ID,
	}
	_, err := s.region.doPut(ServiceInstance, "/v2/volume/"+diskId+"?attach", params)
	return err
}

// ChangeConfig implements cloudprovider.ICloudVM.
func (s *SInstance) ChangeConfig(ctx context.Context, config *cloudprovider.SManagedVMChangeConfig) error {
	panic("unimplemented")
}

// ConvertPublicIpToEip implements cloudprovider.ICloudVM.
// Subtle: this method shadows the method (SInstanceBase).ConvertPublicIpToEip of SInstance.SInstanceBase.
func (s *SInstance) ConvertPublicIpToEip() error {
	panic("unimplemented")
}

// DeleteVM implements cloudprovider.ICloudVM.
func (s *SInstance) DeleteVM(ctx context.Context) error {
	var params = map[string]interface{}{
		"relatedReleaseFlag":    true,
		"deleteCdsSnapshotFlag": true,
		"bccRecycleFlag":        false,
		// "deleteRelatedEnisFlag": true,
		"instanceIds": []string{s.ID},
	}
	_, err := s.region.doPost(ServiceInstance, "/v2/instance/batchDelete", params)
	return err
}

// DeployVM implements cloudprovider.ICloudVM.
func (s *SInstance) DeployVM(ctx context.Context, opts *cloudprovider.SInstanceDeployOptions) error {
	password, err := Aes128EncryptUseSecreteKey(s.region.client.accessKeySecret, opts.Password)
	if err != nil {
		return nil
	}
	var params = map[string]interface{}{
		"adminPass": password,
	}
	_, err = s.region.doPut(ServiceInstance, "/v2/instance/"+s.ID+"?changePass", params)
	return err
}

// DetachDisk implements cloudprovider.ICloudVM.
func (s *SInstance) DetachDisk(ctx context.Context, diskId string) error {
	var params = map[string]interface{}{
		"instanceId": s.ID,
	}
	_, err := s.region.doPut(ServiceInstance, "/v2/volume/"+diskId+"?detach", params)
	return err
}

// GetBios implements cloudprovider.ICloudVM.
func (s *SInstance) GetBios() cloudprovider.TBiosType {
	return ""
}

// GetBootOrder implements cloudprovider.ICloudVM.
func (s *SInstance) GetBootOrder() string {
	return "dcn"
}

// GetCpuSockets implements cloudprovider.ICloudVM.
// Subtle: this method shadows the method (SInstanceBase).GetCpuSockets of SInstance.SInstanceBase.
func (s *SInstance) GetCpuSockets() int {
	return int(s.CPUCount)
}

// GetDescription implements cloudprovider.ICloudVM.
// Subtle: this method shadows the method (SInstanceBase).GetDescription of SInstance.SInstanceBase.
func (s *SInstance) GetDescription() string {
	return s.Desc
}

// GetExpiredAt implements cloudprovider.ICloudVM.
// Subtle: this method shadows the method (SBillingBase).GetExpiredAt of SInstance.SBillingBase.
func (s *SInstance) GetExpiredAt() time.Time {
	return s.ExpireTime
}

// GetFullOsName implements cloudprovider.ICloudVM.
func (s *SInstance) GetFullOsName() string {
	if s.image == nil {
		if err := s.getOsInfo(); err != nil {
			return ""
		}
	}
	return s.image.GetFullOsName()
}

// GetGlobalId implements cloudprovider.ICloudVM.
func (s *SInstance) GetGlobalId() string {
	return s.ID
}

// GetHostname implements cloudprovider.ICloudVM.
func (s *SInstance) GetHostname() string {
	return s.Hostname
}

// GetHypervisor implements cloudprovider.ICloudVM.
func (s *SInstance) GetHypervisor() string {
	return CLOUD_PROVIDER_BAIDU
}

// GetIDisks implements cloudprovider.ICloudVM.
func (s *SInstance) GetIDisks() ([]cloudprovider.ICloudDisk, error) {
	var ret []cloudprovider.ICloudDisk
	disk, err := s.GetDisks()
	if err != nil {
		return nil, err
	}
	for i := range disk {
		disk[i].storage = &SStorage{storageType: disk[i].StorageType, zone: s.host.zone}
		ret = append(ret, &disk[i])
	}
	return ret, nil
}

func (s *SInstance) GetTags() (map[string]string, error) {
	var tag = make(map[string]string)
	for _, val := range s.Tags {
		tag[strings.ReplaceAll(val.TagKey, "/", ":")[6:]] = val.TagValue
	}
	return tag, nil
}

type volumesResp struct {
	NextMarker  string  `json:"nextMarker"`
	Marker      string  `json:"marker"`
	MaxKeys     int     `json:"maxKeys"`
	IsTruncated bool    `json:"isTruncated"`
	Disk        []SDisk `json:"volumes"`
}

// GetIDisks implements cloudprovider.ICloudVM.
func (s *SInstance) GetDisks() ([]SDisk, error) {
	var query = map[string]string{
		"instanceId": s.ID,
		"maxKeys":    "100",
	}
	var marker string
	var ret []SDisk
	for {
		if len(marker) > 0 {
			query["marker"] = marker
		}
		res, err := s.region.doList(ServiceDisk, "/v2/volume", query)
		if err != nil {
			return nil, err
		}
		var r volumesResp
		if err = res.Unmarshal(&r); err != nil {
			return nil, err
		}
		ret = append(ret, r.Disk...)
		if !r.IsTruncated {
			break
		}
		marker = r.NextMarker
	}
	return ret, nil
}

// GetIEIP implements cloudprovider.ICloudVM.
func (s *SInstance) GetIEIP() (cloudprovider.ICloudEIP, error) {
	return nil, nil
}

// GetIHost implements cloudprovider.ICloudVM.
func (s *SInstance) GetIHost() cloudprovider.ICloudHost {
	return s.host
}

// GetIHostId implements cloudprovider.ICloudVM.
// Subtle: this method shadows the method (SInstanceBase).GetIHostId of SInstance.SInstanceBase.
func (s *SInstance) GetIHostId() string {
	return fmt.Sprintf("-%s", s.ZoneName)
}

// GetINics implements cloudprovider.ICloudVM.
func (s *SInstance) GetINics() ([]cloudprovider.ICloudNic, error) {
	return []cloudprovider.ICloudNic{&s.NicInfo}, nil
}

// GetId implements cloudprovider.ICloudVM.
func (s *SInstance) GetId() string {
	return s.ID
}

// GetInstanceType implements cloudprovider.ICloudVM.
func (s *SInstance) GetInstanceType() string {
	return s.Spec
}

// GetInternetMaxBandwidthOut implements cloudprovider.ICloudVM.
// Subtle: this method shadows the method (SInstanceBase).GetInternetMaxBandwidthOut of SInstance.SInstanceBase.
func (s *SInstance) GetInternetMaxBandwidthOut() int {
	return 0
}

// GetMachine implements cloudprovider.ICloudVM.
func (s *SInstance) GetMachine() string {
	return "pc"
}

// GetName implements cloudprovider.ICloudVM.
func (s *SInstance) GetName() string {
	return s.Name
}

// GetOsArch implements cloudprovider.ICloudVM.
func (s *SInstance) GetOsArch() string {
	if s.image == nil {
		if err := s.getOsInfo(); err != nil {
			return ""
		}
	}
	return s.image.OsArch
}

func (s *SInstance) getOsInfo() error {
	var image SImage
	err := s.region.doGet(ServiceImage, "/v2/image/"+s.ImageID, nil, &image)
	if err != nil {
		return err
	}
	s.image = &image
	return nil
}

// GetOsDist implements cloudprovider.ICloudVM.
func (s *SInstance) GetOsDist() string {
	if s.image == nil {
		if err := s.getOsInfo(); err != nil {
			return ""
		}
	}
	return s.image.OsName
}

// GetOsLang implements cloudprovider.ICloudVM.
func (s *SInstance) GetOsLang() string {
	if s.image == nil {
		if err := s.getOsInfo(); err != nil {
			return ""
		}
	}
	return s.image.OsLang
}

// GetOsType implements cloudprovider.ICloudVM.
func (s *SInstance) GetOsType() cloudprovider.TOsType {
	if s.image == nil {
		if err := s.getOsInfo(); err != nil {
			return ""
		}
	}
	return cloudprovider.TOsType(s.image.OsType)
}

// GetOsVersion implements cloudprovider.ICloudVM.
func (s *SInstance) GetOsVersion() string {
	if s.image == nil {
		if err := s.getOsInfo(); err != nil {
			return ""
		}
	}
	return s.image.OsVersion
}

// GetProjectId implements cloudprovider.ICloudVM.
func (s *SInstance) GetProjectId() string {
	return ""
}

// GetSecurityGroupIds implements cloudprovider.ICloudVM.
func (s *SInstance) GetSecurityGroupIds() ([]string, error) {
	return s.NicInfo.SecurityGroups, nil
}

// GetStatus implements cloudprovider.ICloudVM.
func (s *SInstance) GetStatus() string {
	switch s.Status {
	case "Running":
		return api.VM_RUNNING
	case "Stopped":
		return api.VM_READY
	case "Starting":
		return api.VM_STARTING
	case "Stopping":
		return api.VM_STOPPING
	case "Error":
		return api.VM_DEPLOY_FAILED
	default:
		return api.VM_UNKNOWN
	}
}

// GetVNCInfo implements cloudprovider.ICloudVM.
func (s *SInstance) GetVNCInfo(input *cloudprovider.ServerVncInput) (*cloudprovider.ServerVncOutput, error) {
	return nil, nil
}

// GetVcpuCount implements cloudprovider.ICloudVM.
func (s *SInstance) GetVcpuCount() int {
	return s.CPUCount
}

// GetVdi implements cloudprovider.ICloudVM.
func (s *SInstance) GetVdi() string {
	return "vnc"
}

// GetVga implements cloudprovider.ICloudVM.
func (s *SInstance) GetVga() string {
	return "std"
}

// GetVmemSizeMB implements cloudprovider.ICloudVM.
func (s *SInstance) GetVmemSizeMB() int {
	return s.MemoryCapacityInGB * 1024
}

// RebuildRoot implements cloudprovider.ICloudVM.
func (s *SInstance) RebuildRoot(ctx context.Context, config *cloudprovider.SManagedVMRebuildRootConfig) (string, error) {

	var params = map[string]interface{}{
		"imageId": config.ImageId,
		// "adminPass":"adminPass",
	}
	if config.Password != "" {
		password, err := Aes128EncryptUseSecreteKey(s.host.zone.region.client.accessKeySecret, config.Password)
		if err != nil {
			return "", nil
		}
		params["adminPass"] = password
	}
	_, err := s.host.zone.region.doPut(ServiceInstance, "/v2/instance/"+s.ID+"?rebuild", params)
	return "", err
}

// SetSecurityGroups implements cloudprovider.ICloudVM.
func (s *SInstance) SetSecurityGroups(secgroupIds []string) error {
	panic("unimplemented")
}

// StartVM implements cloudprovider.ICloudVM.
func (s *SInstance) StartVM(ctx context.Context) error {
	_, err := s.region.doPut(ServiceInstance, "/v2/instance/"+s.ID+"?start", nil)
	return err
}

// StopVM implements cloudprovider.ICloudVM.
func (s *SInstance) StopVM(ctx context.Context, opts *cloudprovider.ServerStopOptions) error {
	var params = map[string]interface{}{
		"forceStop":        true,
		"stopWithNoCharge": false,
	}
	_, err := s.region.doPut(ServiceInstance, "/v2/instance/"+s.ID+"?stop", params)
	return err
}

// UpdateUserData implements cloudprovider.ICloudVM.
func (s *SInstance) UpdateUserData(userData string) error {
	panic("unimplemented")
}

// UpdateVM implements cloudprovider.ICloudVM.
func (s *SInstance) UpdateVM(ctx context.Context, input cloudprovider.SInstanceUpdateOptions) error {
	panic("unimplemented")
}

var _ cloudprovider.ICloudVM = (*SInstance)(nil)
