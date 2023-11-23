// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package qcloud

import (
	"time"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/utils"
)

// "time"

// {"CpuCoreCount":1,"EniQuantity":1,"GPUAmount":0,"GPUSpec":"","InstanceTypeFamily":"ecs.t1","InstanceTypeId":"ecs.t1.xsmall","LocalStorageCategory":"","MemorySize":0.500000}
// InstanceBandwidthRx":26214400,"InstanceBandwidthTx":26214400,"InstancePpsRx":4500000,"InstancePpsTx":4500000

// Verify that *SInstanceType implements ICloudSku interface
var _ cloudprovider.ICloudSku = (*SInstanceType)(nil)

type SInstanceType struct {
	Zone              string //	可用区。
	InstanceType      string //	实例机型。
	InstanceFamily    string //	实例机型系列。
	GPU               int    //	GPU核数，单位：核。
	CPU               int    //	CPU核数，单位：核。
	Memory            int    //	内存容量，单位：GB。
	CbsSupport        string //	是否支持云硬盘。取值范围：TRUE：表示支持云硬盘；FALSE：表示不支持云硬盘。
	InstanceTypeState string //	机型状态。取值范围：AVAILABLE：表示机型可用；UNAVAILABLE：表示机型不可用。
}

func (self *SInstanceType) GetId() string {
	return self.InstanceType
}

func (self *SInstanceType) GetName() string {
	return self.InstanceType
}

func (self *SInstanceType) GetGlobalId() string {
	return self.InstanceType
}

func (self *SInstanceType) GetCreatedAt() time.Time {

	return time.Now()
}

func (self *SInstanceType) GetDescription() string {
	return ""
}

func (self *SInstanceType) GetStatus() string {
	return ""
}

func (self *SInstanceType) Refresh() error {
	return nil
}

func (self *SInstanceType) IsEmulated() bool {
	return false
}

func (self *SInstanceType) GetSysTags() map[string]string {
	return nil
}

func (self *SInstanceType) GetTags() (map[string]string, error) {
	return nil, nil
}

func (self *SInstanceType) SetTags(tags map[string]string, replace bool) error {
	return nil
}

func (self *SInstanceType) GetZoneID() string {
	return self.Zone
}

func (self *SInstanceType) GetInstanceTypeFamily() string {
	return self.InstanceFamily
}

func (self *SInstanceType) GetInstanceTypeCategory() string {
	return ""
}

func (self *SInstanceType) GetPrepaidStatus() string {
	return ""
}

func (self *SInstanceType) GetPostpaidStatus() string {
	return ""
}

func (self *SInstanceType) GetCpuArch() string {
	return ""
}

func (self *SInstanceType) GetCpuCoreCount() int {
	return self.CPU
}

func (self *SInstanceType) GetMemorySizeMB() int {
	return int(self.Memory * 1024)
}

func (self *SInstanceType) GetOsName() string {
	return ""
}

func (self *SInstanceType) GetSysDiskResizable() bool {
	return false
}

func (self *SInstanceType) GetSysDiskType() string {
	return ""
}

func (self *SInstanceType) GetSysDiskMinSizeGB() int {
	return 0
}

func (self *SInstanceType) GetSysDiskMaxSizeGB() int {
	return 0
}

func (self *SInstanceType) GetAttachedDiskType() string {
	return ""
}

func (self *SInstanceType) GetAttachedDiskSizeGB() int {
	return 0
}

func (self *SInstanceType) GetAttachedDiskCount() int {
	return 0
}

func (self *SInstanceType) GetDataDiskTypes() string {
	return ""
}

func (self *SInstanceType) GetDataDiskMaxCount() int {
	return 0
}

func (self *SInstanceType) GetNicType() string {
	return ""
}

func (self *SInstanceType) GetNicMaxCount() int {
	return 0
}

func (self *SInstanceType) GetGpuAttachable() bool {
	return false
}

func (self *SInstanceType) GetGpuSpec() string {
	return ""
}

func (self *SInstanceType) GetGpuCount() string {
	return ""
}

func (self *SInstanceType) GetGpuMaxCount() int {
	return 0
}

func (self *SInstanceType) Delete() error {
	return nil
}

func (self *SRegion) GetInstanceTypes() ([]SInstanceType, error) {
	params := make(map[string]string)
	params["Region"] = self.Region

	//body, err := self.cvmRequest("DescribeInstanceTypeConfigs", params, true)
	body, err := self.cvmRequest("DescribeInstanceConfigInfos", params, true)
	if err != nil {
		log.Errorf("DescribeInstanceTypeConfigs fail %s", err)
		return nil, err
	}

	instanceTypes := make([]SInstanceType, 0)
	err = body.Unmarshal(&instanceTypes, "InstanceTypeConfigSet")
	if err != nil {
		log.Errorf("Unmarshal instance type details fail %s", err)
		return nil, err
	}
	return instanceTypes, nil
}

func (self *SInstanceType) memoryMB() int {
	return int(self.Memory * 1024)
}

type SLocalDiskType struct {
	Type          string
	PartitionType string
	MinSize       int
	MaxSize       int
}

type SStorageBlockAttr struct {
	Type    string
	MinSize int
	MaxSize int
}

type SExternal struct {
	ReleaseAddress    string
	UnsupportNetworks []string
	StorageBlockAttr  SStorageBlockAttr
}

type SZoneInstanceType struct {
	Zone               string
	InstanceType       string
	InstanceChargeType string
	NetworkCard        int
	Externals          SExternal
	Cpu                int
	Memory             int
	InstanceFamily     string
	TypeName           string
	LocalDiskTypeList  []SLocalDiskType
	Status             string
}

func (self *SRegion) GetZoneInstanceTypes(zoneId string) ([]SZoneInstanceType, error) {
	params := map[string]string{}
	params["Region"] = self.Region
	params["Filters.0.Name"] = "zone"
	params["Filters.0.Values.0"] = zoneId
	body, err := self.cvmRequest("DescribeZoneInstanceConfigInfos", params, true)
	if err != nil {
		return nil, errors.Wrap(err, "DescribeZoneInstanceConfigInfos")
	}
	instanceTypes := []SZoneInstanceType{}
	err = body.Unmarshal(&instanceTypes, "InstanceTypeQuotaSet")
	if err != nil {
		return nil, errors.Wrap(err, "body.Unmarshal")
	}
	return instanceTypes, nil
}

func (self *SRegion) GetZoneLocalStorages(zoneId string) ([]string, error) {
	instanceTypes, err := self.GetZoneInstanceTypes(zoneId)
	if err != nil {
		return nil, errors.Wrap(err, "GetZoneInstanceTypes")
	}
	storages := []string{}
	for _, instanceType := range instanceTypes {
		storage := instanceType.Externals.StorageBlockAttr.Type
		if len(storage) > 0 && !utils.IsInStringArray(storage, storages) {
			storages = append(storages, storage)
		}
		for _, localstorage := range instanceType.LocalDiskTypeList {
			if len(localstorage.Type) > 0 && !utils.IsInStringArray(localstorage.Type, storages) {
				storages = append(storages, localstorage.Type)
			}
		}
	}
	return storages, nil
}

func (self *SRegion) GetISkus() ([]cloudprovider.ICloudSku, error) {
	skus, err := self.GetInstanceTypes()
	if err != nil {
		return nil, errors.Wrapf(err, "GetInstanceTypes")
	}
	ret := []cloudprovider.ICloudSku{}
	for i := range skus {
		ret = append(ret, &skus[i])
	}
	return ret, nil
}
