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
	"math"
	"strconv"
	"time"
	api "yunion.io/x/cloudmux/pkg/apis/compute"

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

// DescribeInstanceTypeAvailable
// 查询指定 instanceType 在zone 下是否售卖
func (self *SRegion) DescribeInstanceTypeAvailable(instanceType, zone string) (bool, error) {
	params := map[string]string{
		"LocationType": "availability-zone",
	}

	params["Filter.1.Name"] = "location"
	params["Filter.1.Value.1"] = zone
	params["Filter.2.Name"] = "instance-type"
	params["Filter.2.Value.1"] = instanceType

	ret := struct {
		InstanceTypeSet []InstanceTypeOffering `xml:"instanceTypeOfferingSet>item"`
		NextToken       string                 `xml:"nextToken"`
	}{}
	err := self.ec2Request("DescribeInstanceTypeOfferings", params, &ret)
	if err != nil {
		return false, err
	}
	return len(ret.InstanceTypeSet) > 0, nil
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

func (self *SRegion) ListPriceLists() (*SInstanceType, error) {
	filters := map[string]string{
		"regionCode": self.RegionId,
		//"operatingSystem": "Linux",
		//"licenseModel":    "No License required",
		//"productFamily":   "Compute Instance",
		//"operation":       "RunInstances",
		//"preInstalledSw":  "NA",
		//"tenancy":         "Shared",
		//"capacitystatus":  "Used",
		//"instanceType":    instanceType,
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
	if len(ret) > 0 {
		return &ret[0], nil
	} else {
		return nil, nil
	}
}

// Price

func (self *SRegion) GetInstanceTypePrice(instanceType string) (*SInstanceType, error) {
	filters := map[string]string{
		"regionCode":      self.RegionId,
		"operatingSystem": "Linux",
		"licenseModel":    "No License required",
		"productFamily":   "Compute Instance",
		"operation":       "RunInstances",
		"preInstalledSw":  "NA",
		"tenancy":         "Shared",
		"capacitystatus":  "Used",
		"instanceType":    instanceType,
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
	if len(ret) > 0 {
		return &ret[0], nil
	} else {
		return nil, errors.Errorf("No Result")
	}
}

type SpotPrice struct {

	// The Availability Zone.
	AvailabilityZone string `xml:"availabilityZone" type:"string"`

	// The instance type.
	InstanceType string `xml:"instanceType" type:"string" enum:"InstanceType"`

	// A general description of the AMI.
	ProductDescription string `xml:"productDescription" type:"string" enum:"RIProductDescription"`

	// The maximum price per unit hour that you are willing to pay for a Spot Instance.
	// We do not recommend using this parameter because it can lead to increased
	// interruptions. If you do not specify this parameter, you will pay the current
	// Spot price.
	//
	// If you specify a maximum price, your instances will be interrupted more frequently
	// than if you do not specify this parameter.
	SpotPrice string `xml:"spotPrice" type:"string"`

	// The date and time the request was created, in UTC format (for example, YYYY-MM-DDTHH:MM:SSZ).
	Timestamp time.Time `xml:"timestamp" type:"timestamp"`
	// contains filtered or unexported fields
}

func (self *SRegion) DescribeSpotPriceHistory(zone, instanceType string) ([]*SpotPrice, error) {

	currentTime := time.Now().UTC()
	startTime := currentTime.Add(-24 * 60 * 60 * time.Second)
	params := map[string]string{
		"AvailabilityZone": zone,
		"StartTime":        startTime.Format("2006-01-02T15:04:05Z"),
		"EndTime":          currentTime.Format("2006-01-02T15:04:05Z"),
	}

	params["ProductDescription.1"] = "Linux/UNIX"
	params["InstanceType.1"] = instanceType

	ret := struct {
		SpotPriceHistorySet []*SpotPrice `xml:"spotPriceHistorySet>item"`
		NextToken           string       `xml:"nextToken"`
	}{}
	err := self.ec2Request("DescribeSpotPriceHistory", params, &ret)
	if err != nil {
		return nil, err
	}
	return ret.SpotPriceHistorySet, nil
}

func (self *SRegion) GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) {
	spotPriceHistory, err := self.DescribeSpotPriceHistory(zoneID, instanceType)
	if err != nil {
		return 0, err
	}
	if len(spotPriceHistory) > 0 {
		newSpotPrice, err := strconv.ParseFloat(spotPriceHistory[0].SpotPrice, 64)
		if err != nil {
			return 0, err
		}
		return newSpotPrice, nil
	}
	return 0, nil
}

func (self *SRegion) GetPostPaidPrice(zoneID, instanceType string) (float64, error) {
	price, err := self.GetInstanceTypePrice(instanceType)
	if err != nil {
		return 0, errors.Wrapf(err, "GetInstanceTypePrice")
	}
	var value float64
	for _, term := range price.Terms.OnDemand {
		for _, dimension := range term.PriceDimensions {
			value = dimension.PricePerUnit.Usd
		}
	}
	return value, nil
}

// GetPrePaidPrice
// PurchaseOption: No Upfront  All Upfront  Partial Upfront
/**
在AWS（Amazon Web Services）中，包年包月的价格模型有三种支付类型，它们分别是：

No Upfront (无预付费)： 在这种支付类型下，无需提前支付任何费用，而是在每个计费周期结束时按月支付使用费用。这是一种相对灵活的选项，适合那些不愿意一次性支付较大费用的用户。

All Upfront (全部预付费)： 在这种支付类型下，需要在合同开始时一次性支付整个合同期的费用，这样您就能够获得更大的折扣。虽然一开始支付较多费用，但在合同期内不需要再支付额外费用。

Partial Upfront (部分预付费)： 这是一种折中的选择，您需要在合同开始时支付部分费用，而余下的费用则按月支付。这样可以降低一次性支付的压力，并且仍然能够获得一定的折扣。
*/
func (self *SRegion) GetPrePaidPrice(zoneID, instanceType string) (float64, error) {
	price, err := self.GetInstanceTypePrice(instanceType)
	if err != nil {
		return 0, errors.Wrapf(err, "GetInstanceTypePrice")
	}
	var value float64

	for _, term := range price.Terms.Reserved {
		attributes := term.TermAttributes
		if attributes.LeaseContractLength == "1yr" && attributes.OfferingClass == "standard" && attributes.PurchaseOption == "All Upfront" {
			for _, dimension := range term.PriceDimensions {
				if dimension.Unit == "Quantity" {
					// 1yr ==> 1month
					value = math.Round(dimension.PricePerUnit.Usd/12*1000) / 1000
				}
			}
		}
	}
	return value, nil
}

func (self *SRegion) GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) {
	available, err := self.DescribeInstanceTypeAvailable(instanceType, zoneID)
	if err != nil {
		return api.SkuStatusSoldout, errors.Wrapf(err, "DescribeInstanceTypeAvailable")
	}
	if available {
		return api.SkuStatusAvailable, nil
	}
	return api.SkuStatusSoldout, nil
}

func (self *SRegion) GetPostPaidStatus(zoneID, instanceType string) (string, error) {
	available, err := self.DescribeInstanceTypeAvailable(instanceType, zoneID)
	if err != nil {
		return api.SkuStatusSoldout, errors.Wrapf(err, "DescribeInstanceTypeAvailable")
	}
	if available {
		return api.SkuStatusAvailable, nil
	}
	return api.SkuStatusSoldout, nil
}

func (self *SRegion) GetPrePaidStatus(zoneID, instanceType string) (string, error) {
	available, err := self.DescribeInstanceTypeAvailable(instanceType, zoneID)
	if err != nil {
		return api.SkuStatusSoldout, errors.Wrapf(err, "DescribeInstanceTypeAvailable")
	}
	if available {
		return api.SkuStatusAvailable, nil
	}
	return api.SkuStatusSoldout, nil
}
