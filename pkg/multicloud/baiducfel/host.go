package baiducfel

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
)

type SHost struct {
	multicloud.SHostBase
	zone *SZone
}

var _ cloudprovider.ICloudHost = (*SHost)(nil)

type billing struct {
	PaymentTiming string      `json:"paymentTiming,omitempty"`
	Reservation   reservation `json:"reservation,omitempty"`
}

type reservation struct {
	ReservationLength   int    `json:"reservationLength,omitempty"`
	ReservationTimeUnit string `json:"reservationTimeUnit,omitempty"`
}

// CreateVM implements cloudprovider.ICloudHost.
func (s *SHost) CreateVM(opts *cloudprovider.SManagedVMCreateConfig) (cloudprovider.ICloudVM, error) {
	// https://cloud.baidu.com/doc/BCC/s/2k3rau7n4

	var dataDisk []map[string]interface{}
	for _, val := range opts.DataDisks {
		dd := map[string]interface{}{
			"storageType": val.StorageType,
			"cdsSizeInGB": val.SizeGB,
		}
		dataDisk = append(dataDisk, dd)
	}
	var tags []BaiduTags
	for k, v := range opts.Tags {
		if len(k) > 0 && len(k) < 65 && len(v) > 0 {
			tags = append(tags, BaiduTags{TagKey: strings.ReplaceAll(k, ":", "/"), TagValue: v})
		}
	}
	var billing = billing{
		PaymentTiming: "Postpaid", // 预支付（Prepaid）和后支付（Postpaid）
	}
	if bc := opts.BillingCycle; bc != nil { // 预付费 prepaid
		if bc.Unit == "W" {
			return nil, errors.New("only support Month and duration in [1,2,3,4,5,6,7,8,9,12,24,36]")
		}
		var count = bc.Count
		
		if bc.Unit == "Y" { //转换成月
			count = bc.Count * 12
		}
		if !slices.Contains(duration, count) {
			return nil, errors.New("only support Month and duration in [1,2,3,4,5,6,7,8,9,12,24,36]")
		}
		billing.PaymentTiming = "Prepaid"
		billing.Reservation = reservation{
			ReservationLength:   count,
			ReservationTimeUnit: "Month",
		}
	}
	var password string
	var err error
	if opts.Password != "" {
		password, err = Aes128EncryptUseSecreteKey(s.zone.region.client.accessKeySecret, opts.Password)
		if err != nil {
			return nil, err
		}
	}
	var params = map[string]interface{}{
		"spec":                opts.InstanceType,
		"rootDiskSizeInGb":    opts.SysDisk.SizeGB,
		"rootDiskStorageType": opts.SysDisk.StorageType,
		"createCdsList":       dataDisk,
		"name":                opts.Name,
		"hostname":            opts.Hostname,
		// "autoSeqSuffix": autoSeqSuffix,
		// "isOpenHostnameDomain": isOpenHostnameDomain,
		"imageId":          opts.ExternalImageId,
		"billing":          billing,
		"zoneName":         s.zone.ZoneName,
		"subnetId":         opts.ExternalNetworkId,
		"securityGroupIds": opts.ExternalSecgroupIds,
		"tags":             tags,
		"userData":         opts.UserData,
		// "keypairId": "",
		// "aspId": "aspId",
		// "specId": "specId", //规格族
		// "resGroupId": "resGroupId",
		// "ehcClusterId": "ehcClusterId"
	}
	if password != "" {
		params["adminPass"] = password
	}
	res, err := s.zone.region.doPost(ServiceInstance, "/v2/instanceBySpec", params)
	if err != nil {
		return nil, err
	}
	var ids []string
	if err = res.Unmarshal(&ids, "instanceIds"); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, errors.New("no instances ids")
	}
	var vm = &SInstance{ID: ids[0]}
	return vm, nil
}

// GetAccessIp implements cloudprovider.ICloudHost.
func (s *SHost) GetAccessIp() string {
	panic("unimplemented")
}

// GetAccessMac implements cloudprovider.ICloudHost.
func (s *SHost) GetAccessMac() string {
	panic("unimplemented")
}

// GetCpuArchitecture implements cloudprovider.ICloudHost.
// Subtle: this method shadows the method (SHostBase).GetCpuArchitecture of SHost.SHostBase.
func (s *SHost) GetCpuArchitecture() string {
	panic("unimplemented")
}

// GetCpuCmtbound implements cloudprovider.ICloudHost.
// Subtle: this method shadows the method (SHostBase).GetCpuCmtbound of SHost.SHostBase.
func (s *SHost) GetCpuCmtbound() float32 {
	panic("unimplemented")
}

// GetCpuCount implements cloudprovider.ICloudHost.
func (s *SHost) GetCpuCount() int {
	panic("unimplemented")
}

// GetCpuDesc implements cloudprovider.ICloudHost.
func (s *SHost) GetCpuDesc() string {
	panic("unimplemented")
}

// GetCpuMhz implements cloudprovider.ICloudHost.
func (s *SHost) GetCpuMhz() int {
	panic("unimplemented")
}

// GetEnabled implements cloudprovider.ICloudHost.
func (s *SHost) GetEnabled() bool {
	panic("unimplemented")
}

// GetGlobalId implements cloudprovider.ICloudHost.
func (s *SHost) GetGlobalId() string {
	return fmt.Sprintf("-%s", s.zone.GetGlobalId())
}

// GetHostStatus implements cloudprovider.ICloudHost.
func (s *SHost) GetHostStatus() string {
	panic("unimplemented")
}

// GetHostType implements cloudprovider.ICloudHost.
func (s *SHost) GetHostType() string {
	panic("unimplemented")
}

// GetIHostNics implements cloudprovider.ICloudHost.
func (s *SHost) GetIHostNics() ([]cloudprovider.ICloudHostNetInterface, error) {
	panic("unimplemented")
}

// GetIStorageById implements cloudprovider.ICloudHost.
func (s *SHost) GetIStorageById(id string) (cloudprovider.ICloudStorage, error) {
	panic("unimplemented")
}

// GetIStorages implements cloudprovider.ICloudHost.
func (s *SHost) GetIStorages() ([]cloudprovider.ICloudStorage, error) {
	panic("unimplemented")
}

// GetIVMById implements cloudprovider.ICloudHost.
func (s *SHost) GetIVMById(id string) (cloudprovider.ICloudVM, error) {
	vm, err := s.zone.region.GetIVMById(id)
	if err != nil {
		return nil, err
	}
	vm.(*SInstance).host = s
	vm.(*SInstance).region = s.zone.region
	return vm, err
}

// GetIVMs implements cloudprovider.ICloudHost.
func (s *SHost) GetIVMs() ([]cloudprovider.ICloudVM, error) {
	panic("unimplemented")
}

// GetId implements cloudprovider.ICloudHost.
func (s *SHost) GetId() string {
	panic("unimplemented")
}

// GetIsMaintenance implements cloudprovider.ICloudHost.
func (s *SHost) GetIsMaintenance() bool {
	panic("unimplemented")
}

// GetMemCmtbound implements cloudprovider.ICloudHost.
// Subtle: this method shadows the method (SHostBase).GetMemCmtbound of SHost.SHostBase.
func (s *SHost) GetMemCmtbound() float32 {
	panic("unimplemented")
}

// GetMemSizeMB implements cloudprovider.ICloudHost.
func (s *SHost) GetMemSizeMB() int {
	panic("unimplemented")
}

// GetName implements cloudprovider.ICloudHost.
func (s *SHost) GetName() string {
	panic("unimplemented")
}

// GetNodeCount implements cloudprovider.ICloudHost.
func (s *SHost) GetNodeCount() int8 {
	panic("unimplemented")
}

// GetOvnVersion implements cloudprovider.ICloudHost.
// Subtle: this method shadows the method (SHostBase).GetOvnVersion of SHost.SHostBase.
func (s *SHost) GetOvnVersion() string {
	panic("unimplemented")
}

// GetReservedMemoryMb implements cloudprovider.ICloudHost.
// Subtle: this method shadows the method (SHostBase).GetReservedMemoryMb of SHost.SHostBase.
func (s *SHost) GetReservedMemoryMb() int {
	panic("unimplemented")
}

// GetSN implements cloudprovider.ICloudHost.
func (s *SHost) GetSN() string {
	panic("unimplemented")
}

// GetStatus implements cloudprovider.ICloudHost.
func (s *SHost) GetStatus() string {
	panic("unimplemented")
}

// GetStorageSizeMB implements cloudprovider.ICloudHost.
func (s *SHost) GetStorageSizeMB() int64 {
	panic("unimplemented")
}

// GetStorageType implements cloudprovider.ICloudHost.
func (s *SHost) GetStorageType() string {
	panic("unimplemented")
}

// GetSysInfo implements cloudprovider.ICloudHost.
func (s *SHost) GetSysInfo() jsonutils.JSONObject {
	panic("unimplemented")
}

// GetVersion implements cloudprovider.ICloudHost.
func (s *SHost) GetVersion() string {
	panic("unimplemented")
}
