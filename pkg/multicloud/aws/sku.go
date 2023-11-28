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

package aws

import (
	"fmt"
	"strconv"
	"time"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

type Product struct {
	ProductFamily string     `json:"productFamily"`
	Attributes    Attributes `json:"attributes"`
	Sku           string     `json:"sku"`
}

type Attributes struct {
	Availabilityzone            string `json:"availabilityzone"`
	Classicnetworkingsupport    string `json:"classicnetworkingsupport"`
	GPUMemory                   string `json:"gpuMemory"`
	Instancesku                 string `json:"instancesku"`
	Marketoption                string `json:"marketoption"`
	RegionCode                  string `json:"regionCode"`
	Vpcnetworkingsupport        string `json:"vpcnetworkingsupport"`
	EnhancedNetworkingSupported string `json:"enhancedNetworkingSupported"`
	IntelTurboAvailable         string `json:"intelTurboAvailable"`
	Memory                      string `json:"memory"`
	DedicatedEbsThroughput      string `json:"dedicatedEbsThroughput"`
	Vcpu                        int    `json:"vcpu"`
	Gpu                         int    `json:"gpu"`
	Capacitystatus              string `json:"capacitystatus"`
	LocationType                string `json:"locationType"`
	Storage                     string `json:"storage"`
	InstanceFamily              string `json:"instanceFamily"`
	OperatingSystem             string `json:"operatingSystem"`
	IntelAvx2Available          string `json:"intelAvx2Available"`
	PhysicalProcessor           string `json:"physicalProcessor"`
	ClockSpeed                  string `json:"clockSpeed"`
	Ecu                         string `json:"ecu"`
	NetworkPerformance          string `json:"networkPerformance"`
	Servicename                 string `json:"servicename"`
	InstanceType                string `json:"instanceType"`
	InstanceSku                 string `json:"instancesku"`
	Tenancy                     string `json:"tenancy"`
	Usagetype                   string `json:"usagetype"`
	NormalizationSizeFactor     string `json:"normalizationSizeFactor"`
	IntelAvxAvailable           string `json:"intelAvxAvailable"`
	ProcessorFeatures           string `json:"processorFeatures"`
	Servicecode                 string `json:"servicecode"`
	LicenseModel                string `json:"licenseModel"`
	CurrentGeneration           string `json:"currentGeneration"`
	PreInstalledSw              string `json:"preInstalledSw"`
	Location                    string `json:"location"`
	ProcessorArchitecture       string `json:"processorArchitecture"`
	Operation                   string `json:"operation"`
	VolumeApiName               string `json:"volumeApiName"`
}

type Terms struct {
	OnDemand map[string]Term `json:"OnDemand"`
	Reserved map[string]Term `json:"Reserved"`
}

type Term struct {
	PriceDimensions map[string]Dimension `json:"priceDimensions"`
	Sku             string               `json:"sku"`
	EffectiveDate   string               `json:"effectiveDate"`
	OfferTermCode   string               `json:"offerTermCode"`
	TermAttributes  TermAttributes       `json:"termAttributes"`
}

type Dimension struct {
	Unit         string       `json:"unit"`
	EndRange     string       `json:"endRange"`
	Description  string       `json:"description"`
	AppliesTo    []string     `json:"appliesTo"`
	RateCode     string       `json:"rateCode"`
	BeginRange   string       `json:"beginRange"`
	PricePerUnit PricePerUnit `json:"pricePerUnit"`
}

type TermAttributes struct {
	LeaseContractLength string `json:"LeaseContractLength"`
	OfferingClass       string `json:"OfferingClass"`
	PurchaseOption      string `json:"PurchaseOption"`
}

type PricePerUnit struct {
	Usd float64 `json:"USD"`
	CNY float64 `json:"CNY"`
}

type SInstanceType struct {
	Product         Product `json:"product"`
	ServiceCode     string  `json:"serviceCode"`
	Terms           Terms   `json:"terms"`
	Version         string  `json:"version"`
	PublicationDate string  `json:"publicationDate"`
}

type InstanceType struct {
	InstanceType string `xml:"instanceType"`
	MemoryInfo   struct {
		SizeInMiB int `xml:"sizeInMiB"`
	} `xml:"memoryInfo"`
}

func (self *SRegion) GetInstanceType(name string) (*InstanceType, error) {
	params := map[string]string{
		"InstanceType.1": name,
	}
	ret := struct {
		InstanceTypeSet []InstanceType `xml:"instanceTypeSet>item"`
		NextToken       string         `xml:"nextToken"`
	}{}
	err := self.ec2Request("DescribeInstanceTypes", params, &ret)
	if err != nil {
		return nil, err
	}
	for i := range ret.InstanceTypeSet {
		if ret.InstanceTypeSet[i].InstanceType == name {
			return &ret.InstanceTypeSet[i], nil
		}
	}
	return nil, errors.Wrapf(cloudprovider.ErrNotFound, name)
}

func (self *SRegion) GetInstanceTypes() ([]SInstanceType, error) {
	filters := map[string]string{
		"regionCode":      self.RegionId,
		"operatingSystem": "Linux",
		"licenseModel":    "No License required",
		"productFamily":   "Compute Instance",
		"operation":       "RunInstances",
		"preInstalledSw":  "NA",
		"tenancy":         "Shared",
		"capacitystatus":  "Used",
	}

	params := []ProductFilter{}

	for k, v := range filters {
		params = append(params, ProductFilter{
			Type:  "TERM_MATCH",
			Field: k,
			Value: v,
		})
	}

	ret := []SInstanceType{}
	var nextToken string
	for {
		parts, _nextToken, err := self.GetProducts("AmazonEC2", params, nextToken)
		if err != nil {
			return nil, err
		}
		ret = append(ret, parts...)
		if len(_nextToken) == 0 || len(parts) == 0 {
			break
		}
		nextToken = _nextToken
	}
	return ret, nil
}

type Sku struct {
	// The instance type. For more information, see Instance types (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html)
	// in the Amazon EC2 User Guide.
	InstanceType string `xml:"instanceType"`

	// cfel 新增字段

	ZoneId string

	// Indicates whether Amazon CloudWatch action based recovery is supported.
	AutoRecoverySupported bool `xml:"autoRecoverySupported"`
	// Indicates whether the instance is a bare metal instance type.
	BareMetal bool `xml:"bareMetal"`
	// Indicates whether the instance type is a burstable performance T instance
	// type. For more information, see Burstable performance instances (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/burstable-performance-instances.html).
	BurstablePerformanceSupported bool `xml:"burstablePerformanceSupported"`
	// Indicates whether the instance type is current generation.
	CurrentGeneration bool `xml:"currentGeneration"`
	// Indicates whether Dedicated Hosts are supported on the instance type.
	DedicatedHostsSupported bool `xml:"dedicatedHostsSupported"`

	// Indicates whether the instance type is eligible for the free tier.
	FreeTierEligible bool `xml:"freeTierEligible" type:"boolean"`

	// Describes the GPU accelerator settings for the instance type.
	GpuInfo GpuInfo `xml:"gpuInfo"`

	// Indicates whether On-Demand hibernation is supported.
	HibernationSupported bool `xml:"hibernationSupported" type:"boolean"`

	// The hypervisor for the instance type.
	Hypervisor string `xml:"hypervisor"`

	// Describes the instance storage for the instance type.
	InstanceStorageInfo InstanceStorageInfo `xml:"instanceStorageInfo" type:"structure"`

	// Indicates whether instance storage is supported.
	InstanceStorageSupported bool `xml:"instanceStorageSupported" type:"boolean"`

	// Describes the memory for the instance type.
	MemoryInfo struct {
		SizeInMiB int64 `xml:"sizeInMiB"`
	} `xml:"memoryInfo" type:"structure"`

	// The supported boot modes. For more information, see Boot modes (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ami-boot.html)
	// in the Amazon EC2 User Guide.
	SupportedBootModes []string `xml:"supportedBootModes>item" locationNameList:"item" type:"list" enum:"BootModeType"`

	// The supported root device types.
	SupportedRootDeviceTypes []string `xml:"supportedRootDeviceTypes>item" locationNameList:"item" type:"list" enum:"RootDeviceType"`

	// Indicates whether the instance type is offered for spot or On-Demand.
	SupportedUsageClasses []string `xml:"supportedUsageClasses>item" locationNameList:"item" type:"list" enum:"UsageClassType"`

	// The supported virtualization types.
	SupportedVirtualizationTypes []string `xml:"supportedVirtualizationTypes>item" locationNameList:"item" type:"list" enum:"VirtualizationType"`

	// Describes the vCPU configurations for the instance type.
	VCpuInfo VCpuInfo `xml:"vCpuInfo"`

	//// Describes the network settings for the instance type.
	//NetworkInfo *NetworkInfo `locationName:"networkInfo" type:"structure"`
	//
	//// Indicates whether Nitro Enclaves is supported.
	//NitroEnclavesSupport *string `locationName:"nitroEnclavesSupport" type:"string" enum:"NitroEnclavesSupport"`
	//
	//// Describes the supported NitroTPM versions for the instance type.
	//NitroTpmInfo *NitroTpmInfo `locationName:"nitroTpmInfo" type:"structure"`
	//
	//// Indicates whether NitroTPM is supported.
	//NitroTpmSupport *string `locationName:"nitroTpmSupport" type:"string" enum:"NitroTpmSupport"`
	//
	//// Describes the placement group settings for the instance type.
	//PlacementGroupInfo *PlacementGroupInfo `locationName:"placementGroupInfo" type:"structure"`
	//
	//// Describes the processor.
	//ProcessorInfo *ProcessorInfo `locationName:"processorInfo" type:"structure"`
	//
	//// Describes the Amazon EBS settings for the instance type.
	//EbsInfo *EbsInfo `locationName:"ebsInfo" type:"structure"`
	//// Describes the FPGA accelerator settings for the instance type.
	//FpgaInfo *FpgaInfo `locationName:"fpgaInfo" type:"structure"`
	//
	//// Describes the Inference accelerator settings for the instance type.
	//InferenceAcceleratorInfo *InferenceAcceleratorInfo `locationName:"inferenceAcceleratorInfo" type:"structure"`
}

func (s Sku) GetId() string {
	return s.InstanceType
}

func (s Sku) GetName() string {
	return s.InstanceType
}

func (s Sku) GetGlobalId() string {
	return s.InstanceType
}

func (s Sku) GetCreatedAt() time.Time {
	return time.Now()
}

func (s Sku) GetDescription() string {
	return ""
}

func (s Sku) GetStatus() string {
	return ""
}

func (s Sku) Refresh() error {
	return nil
}

func (s Sku) IsEmulated() bool {
	return false
}

func (s Sku) GetSysTags() map[string]string {
	return nil
}

func (s Sku) GetTags() (map[string]string, error) {
	return nil, nil
}

func (s Sku) SetTags(tags map[string]string, replace bool) error {
	return nil
}

func (s Sku) GetZoneID() string {
	return s.ZoneId
}

func (s Sku) GetInstanceTypeFamily() string {
	return ""
}

func (s Sku) GetInstanceTypeCategory() string {
	return ""
}

func (s Sku) GetPrepaidStatus() string {
	return ""
}

func (s Sku) GetPostpaidStatus() string {
	return ""
}

func (s Sku) GetCpuArch() string {
	return ""
}

func (s Sku) GetCpuCoreCount() int {
	return int(s.VCpuInfo.DefaultVCpus)
}

func (s Sku) GetMemorySizeMB() int {
	return int(s.MemoryInfo.SizeInMiB)
}

func (s Sku) GetOsName() string {
	return ""
}

func (s Sku) GetSysDiskResizable() bool {
	return false
}

func (s Sku) GetSysDiskType() string {
	return ""
}

func (s Sku) GetSysDiskMinSizeGB() int {
	return int(s.InstanceStorageInfo.TotalSizeInGB * 1024)
}

func (s Sku) GetSysDiskMaxSizeGB() int {
	return int(s.InstanceStorageInfo.TotalSizeInGB * 1024)
}

func (s Sku) GetAttachedDiskType() string {
	return ""
}

func (s Sku) GetAttachedDiskSizeGB() int {
	return 0
}

func (s Sku) GetAttachedDiskCount() int {
	return 0
}

func (s Sku) GetDataDiskTypes() string {
	return ""
}

func (s Sku) GetDataDiskMaxCount() int {
	return 0
}

func (s Sku) GetNicType() string {
	return ""
}

func (s Sku) GetNicMaxCount() int {
	return 0
}

func (s Sku) GetGpuAttachable() bool {
	return s.GpuInfo.TotalGpuMemoryInMiB != 0
}

func (s Sku) GetGpuSpec() string {
	if s.GpuInfo.TotalGpuMemoryInMiB != 0 && len(s.GpuInfo.Gpus) > 0 {
		if s.GpuInfo.Gpus[0].Manufacturer != "" && s.GpuInfo.Gpus[0].Name != "" {
			return fmt.Sprintf("%s %s", s.GpuInfo.Gpus[0].Manufacturer, s.GpuInfo.Gpus[0].Name)
		}
		return ""
	}
	return ""
}

func (s Sku) GetGpuCount() string {
	if s.GpuInfo.TotalGpuMemoryInMiB != 0 && len(s.GpuInfo.Gpus) > 0 {
		return strconv.Itoa(s.GpuInfo.Gpus[0].Count)
	}
	return "0"
}

func (s Sku) GetGpuMaxCount() int {
	if s.GpuInfo.TotalGpuMemoryInMiB != 0 && len(s.GpuInfo.Gpus) > 0 {
		return s.GpuInfo.Gpus[0].Count
	}
	return 0
}

func (s Sku) Delete() error {
	return nil
}

func (s *Sku) GetGPUMemorySizeMB() int {
	return s.GpuInfo.TotalGpuMemoryInMiB
}

func (s *Sku) GetIsBareMetal() bool {
	return s.BareMetal
}

type GpuInfo struct {
	Gpus []struct {
		Count        int    `xml:"count"`
		Manufacturer string `xml:"manufacturer"`
		MemoryInfo   struct {
			SizeInMiB int `xml:"sizeInMiB"`
		} `xml:"memoryInfo"`
		Name string `xml:"name"`
	} `xml:"gpus>item"`
	//Gpus                string `xml:"gpus"`
	TotalGpuMemoryInMiB int `xml:"totalGpuMemoryInMiB"`
}
type VCpuInfo struct {
	// The default number of cores for the instance type.
	DefaultCores int64 `locationName:"defaultCores" type:"integer" xml:"defaultCores"`
	// The default number of threads per core for the instance type.
	DefaultThreadsPerCore int64 `locationName:"defaultThreadsPerCore" type:"integer" xml:"defaultThreadsPerCore"`
	// The default number of vCPUs for the instance type.
	DefaultVCpus int64 `locationName:"defaultVCpus" type:"integer" xml:"defaultVCpus"`
	// The valid number of cores that can be configured for the instance type.
	ValidCores []int64 `locationName:"validCores" locationNameList:"item" type:"list" xml:"validCores>item"`
	// The valid number of threads per core that can be configured for the instance
	// type.
	ValidThreadsPerCore []int64 `locationName:"validThreadsPerCore" locationNameList:"item" type:"list" xml:"validThreadsPerCore>item"`
	// contains filtered or unexported fields
}

type InstanceStorageInfo struct {

	// Describes the disks that are available for the instance type.
	Disks []struct {
		// The number of disks with this configuration.
		Count int64 `xml:"count" type:"integer"`

		// The size of the disk in GB.
		SizeInGB int64 `xml:"sizeInGB" type:"long"`

		// The type of disk.
		Type string `xml:"type" type:"string" enum:"DiskType"`
		// contains filtered or unexported fields
	} `xml:"disks>item" locationNameList:"item" type:"list"`

	// Indicates whether data is encrypted at rest.
	EncryptionSupport string `xml:"encryptionSupport" type:"string" enum:"InstanceStorageEncryptionSupport"`

	// Indicates whether non-volatile memory express (NVMe) is supported.
	NvmeSupport string `xml:"nvmeSupport" type:"string" enum:"EphemeralNvmeSupport"`

	// The total size of the disks, in GB.
	TotalSizeInGB int64 `xml:"totalSizeInGB" type:"long"`
	// contains filtered or unexported fields
}

func (self *SRegion) DescribeInstanceTypes(arch string, nextToken string) ([]Sku, string, error) {
	params := map[string]string{}
	if len(nextToken) > 0 {
		params["NextToken"] = nextToken
	}
	idx := 1
	if len(arch) > 0 {
		params[fmt.Sprintf("Filter.%d.Name", idx)] = "processor-info.supported-architecture"
		params[fmt.Sprintf("Filter.%d.Value.1", idx)] = arch
		idx++
	}
	ret := struct {
		InstanceTypeSet []Sku  `xml:"instanceTypeSet>item"`
		NextToken       string `xml:"nextToken"`
	}{}
	err := self.ec2Request("DescribeInstanceTypes", params, &ret)
	if err != nil {
		return nil, "", err
	}
	return ret.InstanceTypeSet, ret.NextToken, nil
}

// DescribeInstanceTypesAll 获取某一个 region 下所有 InstanceTypes
func (self *SRegion) DescribeInstanceTypesAll() ([]Sku, error) {
	ret := []Sku{}
	var nextToken string
	for {
		parts, _nextToken, err := self.DescribeInstanceTypes("", nextToken)
		if err != nil {
			return nil, err
		}
		ret = append(ret, parts...)
		if len(_nextToken) == 0 || len(parts) == 0 {
			break
		}
		nextToken = _nextToken
	}

	//for i, sku := range ret {
	//	if sku.GpuInfo.Gpus != nil {
	//		fmt.Println(i, sku)
	//	}
	//}
	return ret, nil
}

type InstanceTypeOffering struct {

	// The instance type. For more information, see Instance types (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html)
	// in the Amazon EC2 User Guide.
	InstanceType string `xml:"instanceType" type:"string" enum:"InstanceType"`

	// The identifier for the location. This depends on the location type. For example,
	// if the location type is region, the location is the Region code (for example,
	// us-east-2.)
	Location string `xml:"location" type:"string"`

	// The location type.
	LocationType string `xml:"locationType" type:"string" enum:"LocationType"`
}

func (self *SRegion) DescribeInstanceTypeOfferings(nextToken string) ([]InstanceTypeOffering, string, error) {
	params := map[string]string{
		"LocationType": "availability-zone",
	}

	if len(nextToken) > 0 {
		params["NextToken"] = nextToken
	}
	//idx := 1
	//if len(arch) > 0 {
	//	params[fmt.Sprintf("Filter.%d.Name", idx)] = "processor-info.supported-architecture"
	//	params[fmt.Sprintf("Filter.%d.Value.1", idx)] = arch
	//	idx++
	//}
	ret := struct {
		InstanceTypeSet []InstanceTypeOffering `xml:"instanceTypeOfferingSet>item"`
		NextToken       string                 `xml:"nextToken"`
	}{}
	err := self.ec2Request("DescribeInstanceTypeOfferings", params, &ret)
	if err != nil {
		return nil, "", err
	}
	return ret.InstanceTypeSet, ret.NextToken, nil
}

// DescribeInstanceTypeOfferingsAll 获取某一个 region 下所有 InstanceTypeOfferings
func (self *SRegion) DescribeInstanceTypeOfferingsAll() ([]InstanceTypeOffering, error) {
	ret := []InstanceTypeOffering{}
	var nextToken string
	for {
		parts, _nextToken, err := self.DescribeInstanceTypeOfferings(nextToken)
		if err != nil {
			return nil, err
		}
		ret = append(ret, parts...)
		if len(_nextToken) == 0 || len(parts) == 0 {
			break
		}
		nextToken = _nextToken
	}

	return ret, nil
}

// GetRegionZoneInstanceType 获取带zone 的 instanceType
func (self *SRegion) GetRegionZoneInstanceType() ([]Sku, error) {

	existMap := make(map[string]*Sku)
	instanceTypesAll, err := self.DescribeInstanceTypesAll()
	if err != nil {
		return nil, errors.Wrapf(err, "DescribeInstanceTypesAll")
	}
	for _, instanceType := range instanceTypesAll {
		_instanceType := instanceType
		existMap[instanceType.InstanceType] = &_instanceType
	}

	skus := make([]Sku, 0)
	offerings, err := self.DescribeInstanceTypeOfferingsAll()
	if err != nil {
		return nil, errors.Wrapf(err, "DescribeInstanceTypeOfferingsAll")
	}
	for _, offering := range offerings {
		sku, ok := existMap[offering.InstanceType]
		if ok {
			sku.ZoneId = offering.Location
			skus = append(skus, *sku)
		}
	}
	return skus, nil
}

func (self *SRegion) GetISkus() ([]cloudprovider.ICloudSku, error) {
	skus, err := self.GetRegionZoneInstanceType()
	if err != nil {
		return nil, errors.Wrapf(err, "GetInstanceTypes")
	}
	ret := []cloudprovider.ICloudSku{}
	for i := range skus {
		ret = append(ret, &skus[i])
	}
	return ret, nil
}
