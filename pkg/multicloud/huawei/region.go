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

package huawei

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/cloudmux/pkg/multicloud/huawei/client"
	"yunion.io/x/cloudmux/pkg/multicloud/huawei/obs"
)

type Locales struct {
	EnUs string `json:"en-us"`
	ZhCN string `json:"zh-cn"`
}

// Verify that *SInstanceType implements ICloudSkuPrice
var _ cloudprovider.ICloudSkuPrice = (*SRegion)(nil)

// https://support.huaweicloud.com/api-iam/zh-cn_topic_0067148043.html
type SRegion struct {
	multicloud.SRegion

	client    *SHuaweiClient
	ecsClient *client.Client
	obsClient *obs.ObsClient // 对象存储client.请勿直接引用。

	Description    string  `json:"description"`
	ID             string  `json:"id"`
	Locales        Locales `json:"locales"`
	ParentRegionID string  `json:"parent_region_id"`
	Type           string  `json:"type"`

	izones []cloudprovider.ICloudZone
	ivpcs  []cloudprovider.ICloudVpc

	storageCache *SStoragecache
}

func (self *SRegion) GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) {
	// huawei 暂不支持 spot 价格 api 查询
	return -1, nil
}

// 创建ECS高频场景示例
// docs (https://support.huaweicloud.com/api-ecs/ecs_04_0007.html#ecs_04_0007__section413314335610)

// GetPostPaidPrice 查询按需产品价格
// https://support.huaweicloud.com/api-oce/bcloud_01001.html
func (self *SRegion) GetPostPaidPrice(zoneID, instanceType string) (float64, error) {

	params := map[string]interface{}{
		"project_id": self.ecsClient.Bills.ProjectId,
		"product_infos": []interface{}{map[string]interface{}{
			"id":                 time.Now().String(),
			"cloud_service_type": "hws.service.type.ec2",
			"resource_type":      "hws.resource.type.vm",
			"resource_spec":      instanceType + ".linux",
			"region":             self.ID,
			"available_zone":     zoneID,
			"usage_factor":       "Duration",
			"usage_value":        1,
			"usage_measure_id":   4,
			"subscription_num":   1,
		}},
	}

	resp, err := self.ecsClient.Bills.PerformActionWithoutProject("on-demand-resources", "", jsonutils.Marshal(params), "")
	if err != nil {
		return -1, err
	}
	getString, err := resp.GetString("amount")
	f, err := strconv.ParseFloat(getString, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

// GetPrePaidPrice 查询包年/包月产品价格
// https://support.huaweicloud.com/api-oce/bcloud_01002.html
func (self *SRegion) GetPrePaidPrice(zoneID, instanceType string) (float64, error) {
	params := map[string]interface{}{
		"project_id": self.ecsClient.Bills.ProjectId,
		"product_infos": []interface{}{map[string]interface{}{
			"id":                 time.Now().String(),
			"cloud_service_type": "hws.service.type.ec2",
			"resource_type":      "hws.resource.type.vm",
			"resource_spec":      instanceType + ".linux",
			"region":             self.ID,
			"available_zone":     zoneID,
			"period_type":        2, // 订购包年/包月产品的周期类型。 0:天2:月3:年4:小时
			"period_num":         1, // 订购包年/包月产品的周期数。
			"subscription_num":   1, // 订购包年/包月产品的数量。
		}},
	}

	resp, err := self.ecsClient.Bills.PerformActionWithoutProject("period-resources/subscribe-rate", "", jsonutils.Marshal(params), "")
	if err != nil {
		return -1, err
	}
	l := new(ListRateOnPeriodDetailResp)
	err = resp.Unmarshal(l)
	if err != nil {
		return -1, err
	}
	return l.OfficialWebsiteRatingResult.OfficialWebsiteAmount, nil
}

type ListRateOnPeriodDetailResp struct {
	OfficialWebsiteRatingResult struct {
		OfficialWebsiteAmount            float64     `json:"official_website_amount"`
		InstallmentOfficialWebsiteAmount interface{} `json:"installment_official_website_amount"`
		InstallmentPeriodType            interface{} `json:"installment_period_type"`
		MeasureID                        int         `json:"measure_id"`
		ProductRatingResults             []struct {
			ID                               string      `json:"id"`
			ProductID                        string      `json:"product_id"`
			OfficialWebsiteAmount            float64     `json:"official_website_amount"`
			MeasureID                        int         `json:"measure_id"`
			InstallmentOfficialWebsiteAmount interface{} `json:"installment_official_website_amount"`
			InstallmentPeriodType            interface{} `json:"installment_period_type"`
		} `json:"product_rating_results"`
	} `json:"official_website_rating_result"`
	OptionalDiscountRatingResults []interface{} `json:"optional_discount_rating_results"`
	Currency                      string        `json:"currency"`
}

// 查询规格资源是否可购买/资源是否售罄
// https://support.huaweicloud.com/api-ecs/ecs_04_0007.html#ecs_04_0007__section413314335610

func (self *SRegion) GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) {
	// TODO 查询规格销售策略
	// https://support.huaweicloud.com/api-ecs/ecs_02_0403.html
	return self.GetInstanceTypeStatus(zoneID, instanceType)
}

func (self *SRegion) GetPostPaidStatus(zoneID, instanceType string) (string, error) {
	return self.GetInstanceTypeStatus(zoneID, instanceType)
}

func (self *SRegion) GetPrePaidStatus(zoneID, instanceType string) (string, error) {
	return self.GetInstanceTypeStatus(zoneID, instanceType)
}

func (self *SRegion) GetClient() *SHuaweiClient {
	return self.client
}

func (self *SRegion) list(service, resource string, query url.Values) (jsonutils.JSONObject, error) {
	return self.client.list(service, self.ID, resource, query)
}

func (self *SRegion) delete(service, resource string) (jsonutils.JSONObject, error) {
	return self.client.delete(service, self.ID, resource)
}

func (self *SRegion) put(service, resource string, params map[string]interface{}) (jsonutils.JSONObject, error) {
	return self.client.put(service, self.ID, resource, params)
}

func (self *SRegion) post(service, resource string, params map[string]interface{}) (jsonutils.JSONObject, error) {
	return self.client.post(service, self.ID, resource, params)
}

func (self *SRegion) patch(service, resource string, query url.Values, params map[string]interface{}) (jsonutils.JSONObject, error) {
	return self.client.patch(service, self.ID, resource, query, params)
}

func (self *SRegion) getECSClient() (*client.Client, error) {
	var err error

	if len(self.client.projectId) > 0 {
		project, err := self.client.GetProjectById(self.client.projectId)
		if err != nil {
			return nil, err
		}

		if !strings.Contains(project.Name, self.ID) {
			return nil, errors.Errorf("region %s and project %s mismatch", self.ID, project.Name)
		}
	}

	if self.ecsClient == nil {
		self.ecsClient, err = self.client.newRegionAPIClient(self.ID)
		if err != nil {
			return nil, err
		}
	}

	return self.ecsClient, err
}

func (self *SRegion) getOBSEndpoint() string {
	return getOBSEndpoint(self.GetId())
}

func (self *SRegion) getOBSClient(signType obs.SignatureType) (*obs.ObsClient, error) {
	if self.obsClient == nil {
		obsClient, err := self.client.getOBSClient(self.GetId(), signType)
		if err != nil {
			return nil, err
		}

		self.obsClient = obsClient
	}

	return self.obsClient, nil
}

func (self *SRegion) fetchZones() error {
	zones := make([]SZone, 0)
	err := doListAll(self.ecsClient.Zones.List, nil, &zones)
	if err != nil {
		return err
	}

	self.izones = make([]cloudprovider.ICloudZone, 0)
	for i := range zones {
		zone := zones[i]
		zone.region = self
		self.izones = append(self.izones, &zone)
	}
	return nil
}

func (self *SRegion) fetchIVpcs() error {
	// https://support.huaweicloud.com/api-vpc/zh-cn_topic_0020090625.html
	vpcs := make([]SVpc, 0)
	querys := map[string]string{
		"limit": "2048",
	}
	err := doListAllWithMarker(self.ecsClient.Vpcs.List, querys, &vpcs)
	if err != nil {
		return err
	}

	self.ivpcs = make([]cloudprovider.ICloudVpc, 0)
	for i := range vpcs {
		vpc := vpcs[i]
		vpc.region = self
		self.ivpcs = append(self.ivpcs, &vpc)
	}
	return nil
}

func (self *SRegion) GetIVMById(id string) (cloudprovider.ICloudVM, error) {
	if len(id) == 0 {
		return nil, errors.Wrap(cloudprovider.ErrNotFound, "SRegion.GetIVMById")
	}

	instance, err := self.GetInstanceByID(id)
	if err != nil {
		return nil, err
	}
	return &instance, err
}

func (self *SRegion) GetIDiskById(id string) (cloudprovider.ICloudDisk, error) {
	return self.GetDisk(id)
}

func (self *SRegion) GetGeographicInfo() cloudprovider.SGeographicInfo {
	if info, ok := LatitudeAndLongitude[self.ID]; ok {
		return info
	}
	return cloudprovider.SGeographicInfo{}
}

func (self *SRegion) GetILoadBalancers() ([]cloudprovider.ICloudLoadbalancer, error) {
	elbs, err := self.GetLoadBalancers()
	if err != nil {
		return nil, err
	}

	ielbs := make([]cloudprovider.ICloudLoadbalancer, len(elbs))
	for i := range elbs {
		elbs[i].region = self
		ielbs[i] = &elbs[i]
	}

	return ielbs, nil
}

func (self *SRegion) GetLoadBalancers() ([]SLoadbalancer, error) {
	lbs := []SLoadbalancer{}
	params := url.Values{}
	return lbs, self.lbListAll("elb/loadbalancers", params, "loadbalancers", &lbs)
}

func (self *SRegion) GetILoadBalancerById(id string) (cloudprovider.ICloudLoadbalancer, error) {
	elb, err := self.GetLoadbalancer(id)
	if err != nil {
		return nil, err
	}
	return elb, nil
}

func (self *SRegion) GetILoadBalancerAclById(aclId string) (cloudprovider.ICloudLoadbalancerAcl, error) {
	acl, err := self.GetLoadBalancerAcl(aclId)
	if err != nil {
		return nil, err
	}
	return acl, nil
}

func (self *SRegion) GetILoadBalancerCertificateById(certId string) (cloudprovider.ICloudLoadbalancerCertificate, error) {
	cert, err := self.GetLoadBalancerCertificate(certId)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func (self *SRegion) CreateILoadBalancerCertificate(cert *cloudprovider.SLoadbalancerCertificate) (cloudprovider.ICloudLoadbalancerCertificate, error) {
	ret, err := self.CreateLoadBalancerCertificate(cert)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (self *SRegion) GetILoadBalancerAcls() ([]cloudprovider.ICloudLoadbalancerAcl, error) {
	ret, err := self.GetLoadBalancerAcls("")
	if err != nil {
		return nil, err
	}

	iret := make([]cloudprovider.ICloudLoadbalancerAcl, len(ret))
	for i := range ret {
		ret[i].region = self
		iret[i] = &ret[i]
	}
	return iret, nil
}

func (self *SRegion) GetILoadBalancerCertificates() ([]cloudprovider.ICloudLoadbalancerCertificate, error) {
	ret, err := self.GetLoadBalancerCertificates()
	if err != nil {
		return nil, err
	}

	iret := make([]cloudprovider.ICloudLoadbalancerCertificate, len(ret))
	for i := range ret {
		ret[i].region = self
		iret[i] = &ret[i]
	}
	return iret, nil
}

// https://support.huaweicloud.com/api-iam/zh-cn_topic_0057845622.html
func (self *SRegion) GetId() string {
	return self.ID
}

func (self *SRegion) GetName() string {
	return fmt.Sprintf("%s %s", CLOUD_PROVIDER_HUAWEI_CN, self.Locales.ZhCN)
}

func (self *SRegion) GetI18n() cloudprovider.SModelI18nTable {
	en := fmt.Sprintf("%s %s", CLOUD_PROVIDER_HUAWEI_EN, self.Locales.EnUs)
	table := cloudprovider.SModelI18nTable{}
	table["name"] = cloudprovider.NewSModelI18nEntry(self.GetName()).CN(self.GetName()).EN(en)
	return table
}

func (self *SRegion) GetGlobalId() string {
	return fmt.Sprintf("%s/%s", self.client.GetAccessEnv(), self.ID)
}

func (self *SRegion) GetStatus() string {
	return api.CLOUD_REGION_STATUS_INSERVER
}

func (self *SRegion) Refresh() error {
	return nil
}

func (self *SRegion) IsEmulated() bool {
	return false
}

func (self *SRegion) GetLatitude() float32 {
	if locationInfo, ok := LatitudeAndLongitude[self.ID]; ok {
		return locationInfo.Latitude
	}
	return 0.0
}

func (self *SRegion) GetLongitude() float32 {
	if locationInfo, ok := LatitudeAndLongitude[self.ID]; ok {
		return locationInfo.Longitude
	}
	return 0.0
}

func (self *SRegion) fetchInfrastructure() error {
	_, err := self.getECSClient()
	if err != nil {
		return err
	}

	if err := self.fetchZones(); err != nil {
		return err
	}

	if err := self.fetchIVpcs(); err != nil {
		return err
	}

	for i := 0; i < len(self.ivpcs); i += 1 {
		vpc := self.ivpcs[i].(*SVpc)
		wire := SWire{region: self, vpc: vpc}
		vpc.addWire(&wire)

		for j := 0; j < len(self.izones); j += 1 {
			zone := self.izones[j].(*SZone)
			zone.addWire(&wire)
		}
	}
	return nil
}

func (self *SRegion) GetIZones() ([]cloudprovider.ICloudZone, error) {
	if self.izones == nil {
		var err error
		err = self.fetchInfrastructure()
		if err != nil {
			return nil, err
		}
	}
	return self.izones, nil
}

func (self *SRegion) GetIVpcs() ([]cloudprovider.ICloudVpc, error) {
	if self.ivpcs == nil {
		err := self.fetchInfrastructure()
		if err != nil {
			return nil, err
		}
	}
	return self.ivpcs, nil
}

func (self *SRegion) GetIEips() ([]cloudprovider.ICloudEIP, error) {
	eips, err := self.GetEips("", nil)
	if err != nil {
		return nil, err
	}

	ret := []cloudprovider.ICloudEIP{}
	for i := 0; i < len(eips); i += 1 {
		eips[i].region = self
		ret = append(ret, &eips[i])
	}
	return ret, nil
}

func (self *SRegion) GetIVpcById(id string) (cloudprovider.ICloudVpc, error) {
	ivpcs, err := self.GetIVpcs()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(ivpcs); i += 1 {
		if ivpcs[i].GetGlobalId() == id {
			return ivpcs[i], nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SRegion) GetIZoneById(id string) (cloudprovider.ICloudZone, error) {
	izones, err := self.GetIZones()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(izones); i += 1 {
		if izones[i].GetGlobalId() == id {
			return izones[i], nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SRegion) GetIEipById(eipId string) (cloudprovider.ICloudEIP, error) {
	eip, err := self.GetEip(eipId)
	if err != nil {
		return nil, err
	}
	return eip, nil
}

func (self *SRegion) DeleteSecurityGroup(id string) error {
	_, err := self.delete(SERVICE_VPC, "vpc/security-groups/"+id)
	return err
}

func (self *SRegion) GetISecurityGroupById(secgroupId string) (cloudprovider.ICloudSecurityGroup, error) {
	return self.GetSecurityGroup(secgroupId)
}

func (self *SRegion) GetISecurityGroups() ([]cloudprovider.ICloudSecurityGroup, error) {
	groups, err := self.GetSecurityGroups("", "")
	if err != nil {
		return nil, errors.Wrapf(err, "GetSecurityGroups")
	}
	ret := []cloudprovider.ICloudSecurityGroup{}
	for i := range groups {
		groups[i].region = self
		ret = append(ret, &groups[i])
	}
	return ret, nil
}

func (self *SRegion) CreateISecurityGroup(opts *cloudprovider.SecurityGroupCreateInput) (cloudprovider.ICloudSecurityGroup, error) {
	return self.CreateSecurityGroup(opts)
}

// https://support.huaweicloud.com/api-vpc/zh-cn_topic_0020090608.html
func (self *SRegion) CreateIVpc(opts *cloudprovider.VpcCreateOptions) (cloudprovider.ICloudVpc, error) {
	return self.CreateVpc(opts.NAME, opts.CIDR, opts.Desc)
}

func (self *SRegion) CreateVpc(name, cidr, desc string) (*SVpc, error) {
	params := map[string]interface{}{
		"vpc": map[string]string{
			"name":        name,
			"cidr":        cidr,
			"description": desc,
		},
	}
	vpc := &SVpc{region: self}
	return vpc, DoCreate(self.ecsClient.Vpcs.Create, jsonutils.Marshal(params), vpc)
}

// https://support.huaweicloud.com/api-vpc/zh-cn_topic_0020090596.html
// size: 1Mbit/s~2000Mbit/s
// bgpType: 5_telcom，5_union，5_bgp，5_sbgp.
// 东北-大连：5_telcom、5_union
// 华南-广州：5_sbgp
// 华东-上海二：5_sbgp
// 华北-北京一：5_bgp、5_sbgp
// 亚太-香港：5_bgp
func (self *SRegion) CreateEIP(opts *cloudprovider.SEip) (cloudprovider.ICloudEIP, error) {
	eip, err := self.AllocateEIP(opts)
	if err != nil {
		return nil, err
	}
	err = cloudprovider.WaitStatus(eip, api.EIP_STATUS_READY, 5*time.Second, time.Minute)
	return eip, err
}

func (self *SRegion) GetISnapshots() ([]cloudprovider.ICloudSnapshot, error) {
	snapshots, err := self.GetSnapshots("", "")
	if err != nil {
		log.Errorf("self.GetSnapshots fail %s", err)
		return nil, err
	}

	ret := make([]cloudprovider.ICloudSnapshot, len(snapshots))
	for i := 0; i < len(snapshots); i += 1 {
		snapshots[i].region = self
		ret[i] = &snapshots[i]
	}
	return ret, nil
}

func (self *SRegion) GetISnapshotById(snapshotId string) (cloudprovider.ICloudSnapshot, error) {
	snapshot, err := self.GetSnapshotById(snapshotId)
	return &snapshot, err
}

func (self *SRegion) GetIHosts() ([]cloudprovider.ICloudHost, error) {
	iHosts := make([]cloudprovider.ICloudHost, 0)

	izones, err := self.GetIZones()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(izones); i += 1 {
		iZoneHost, err := izones[i].GetIHosts()
		if err != nil {
			return nil, err
		}
		iHosts = append(iHosts, iZoneHost...)
	}
	return iHosts, nil
}

func (self *SRegion) GetIHostById(id string) (cloudprovider.ICloudHost, error) {
	izones, err := self.GetIZones()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(izones); i += 1 {
		ihost, err := izones[i].GetIHostById(id)
		if err == nil {
			return ihost, nil
		} else if errors.Cause(err) != cloudprovider.ErrNotFound {
			return nil, err
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SRegion) GetIStorages() ([]cloudprovider.ICloudStorage, error) {
	iStores := make([]cloudprovider.ICloudStorage, 0)

	izones, err := self.GetIZones()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(izones); i += 1 {
		iZoneStores, err := izones[i].GetIStorages()
		if err != nil {
			return nil, err
		}
		iStores = append(iStores, iZoneStores...)
	}
	return iStores, nil
}

func (self *SRegion) GetIStorageById(id string) (cloudprovider.ICloudStorage, error) {
	izones, err := self.GetIZones()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(izones); i += 1 {
		istore, err := izones[i].GetIStorageById(id)
		if err == nil {
			return istore, nil
		} else if errors.Cause(err) != cloudprovider.ErrNotFound {
			return nil, err
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SRegion) GetProvider() string {
	return CLOUD_PROVIDER_HUAWEI
}

func (self *SRegion) GetCloudEnv() string {
	return self.client.cloudEnv
}

func (self *SRegion) CreateSecurityGroup(opts *cloudprovider.SecurityGroupCreateInput) (*SSecurityGroup, error) {
	params := map[string]interface{}{
		"name":                  opts.Name,
		"description":           opts.Desc,
		"enterprise_project_id": opts.ProjectId,
	}
	resp, err := self.post(SERVICE_VPC, "vpc/security-groups", map[string]interface{}{"security_group": params})
	if err != nil {
		return nil, err
	}
	ret := &SSecurityGroup{region: self}
	return ret, resp.Unmarshal(ret, "security_group")
}

func (self *SRegion) CreateILoadBalancer(loadbalancer *cloudprovider.SLoadbalancerCreateOptions) (cloudprovider.ICloudLoadbalancer, error) {
	ret, err := self.CreateLoadBalancer(loadbalancer)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (self *SRegion) CreateILoadBalancerAcl(acl *cloudprovider.SLoadbalancerAccessControlList) (cloudprovider.ICloudLoadbalancerAcl, error) {
	ret, err := self.CreateLoadBalancerAcl(acl)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (region *SRegion) GetIBuckets() ([]cloudprovider.ICloudBucket, error) {
	iBuckets, err := region.client.getIBuckets()
	if err != nil {
		return nil, errors.Wrap(err, "getIBuckets")
	}
	ret := make([]cloudprovider.ICloudBucket, 0)
	for i := range iBuckets {
		// huawei OBS is shared across projects
		if iBuckets[i].GetLocation() == region.GetId() {
			ret = append(ret, iBuckets[i])
		}
	}
	return ret, nil
}

func str2StorageClass(storageClassStr string) (obs.StorageClassType, error) {
	if strings.EqualFold(storageClassStr, string(obs.StorageClassStandard)) {
		return obs.StorageClassStandard, nil
	} else if strings.EqualFold(storageClassStr, string(obs.StorageClassWarm)) {
		return obs.StorageClassWarm, nil
	} else if strings.EqualFold(storageClassStr, string(obs.StorageClassCold)) {
		return obs.StorageClassCold, nil
	} else {
		return obs.StorageClassStandard, errors.Error("unsupported storageClass")
	}
}

func (region *SRegion) CreateIBucket(name string, storageClassStr string, aclStr string) error {
	obsClient, err := region.getOBSClient("")
	if err != nil {
		return errors.Wrap(err, "region.getOBSClient")
	}
	input := &obs.CreateBucketInput{}
	input.Bucket = name
	input.Location = region.GetId()
	if len(aclStr) > 0 {
		if strings.EqualFold(aclStr, string(obs.AclPrivate)) {
			input.ACL = obs.AclPrivate
		} else if strings.EqualFold(aclStr, string(obs.AclPublicRead)) {
			input.ACL = obs.AclPublicRead
		} else if strings.EqualFold(aclStr, string(obs.AclPublicReadWrite)) {
			input.ACL = obs.AclPublicReadWrite
		} else {
			return errors.Error("unsupported acl")
		}
	}
	if len(storageClassStr) > 0 {
		input.StorageClass, err = str2StorageClass(storageClassStr)
		if err != nil {
			return err
		}
	}
	_, err = obsClient.CreateBucket(input)
	if err != nil {
		return errors.Wrap(err, "obsClient.CreateBucket")
	}
	region.client.invalidateIBuckets()
	return nil
}

func obsHttpCode(err error) int {
	switch httpErr := err.(type) {
	case obs.ObsError:
		return httpErr.StatusCode
	case *obs.ObsError:
		return httpErr.StatusCode
	}
	return -1
}

func (region *SRegion) DeleteIBucket(name string) error {
	obsClient, err := region.getOBSClient("")
	if err != nil {
		return errors.Wrap(err, "region.getOBSClient")
	}
	_, err = obsClient.DeleteBucket(name)
	if err != nil {
		if obsHttpCode(err) == 404 {
			return nil
		}
		log.Debugf("%#v %s", err, err)
		return errors.Wrap(err, "DeleteBucket")
	}
	region.client.invalidateIBuckets()
	return nil
}

func (region *SRegion) HeadBucket(name string) (*obs.BaseModel, error) {
	obsClient, err := region.getOBSClient("")
	if err != nil {
		return nil, errors.Wrap(err, "region.getOBSClient")
	}
	return obsClient.HeadBucket(name)
}

func (region *SRegion) IBucketExist(name string) (bool, error) {
	_, err := region.HeadBucket(name)
	if err != nil {
		if obsHttpCode(err) == 404 {
			return false, nil
		} else {
			return false, errors.Wrap(err, "HeadBucket")
		}
	}
	return true, nil
}

func (region *SRegion) GetIBucketById(name string) (cloudprovider.ICloudBucket, error) {
	return cloudprovider.GetIBucketById(region, name)
}

func (region *SRegion) GetIBucketByName(name string) (cloudprovider.ICloudBucket, error) {
	return region.GetIBucketById(name)
}

func (self *SRegion) GetSkus(zoneId string) ([]cloudprovider.ICloudSku, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SRegion) GetIElasticcaches() ([]cloudprovider.ICloudElasticcache, error) {
	caches, err := self.GetElasticCaches()
	if err != nil {
		return nil, err
	}

	icaches := make([]cloudprovider.ICloudElasticcache, len(caches))
	for i := range caches {
		caches[i].region = self
		icaches[i] = &caches[i]
	}

	return icaches, nil
}

func (region *SRegion) GetCapabilities() []string {
	return region.client.GetCapabilities()
}

func (self *SRegion) GetDiskTypes() ([]SDiskType, error) {
	ret, err := self.ecsClient.Disks.GetDiskTypes()
	if err != nil {
		return nil, errors.Wrap(err, "GetDiskTypes")
	}

	dts := []SDiskType{}
	_ret := jsonutils.NewArray(ret.Data...)
	err = _ret.Unmarshal(&dts)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	return dts, nil
}

func (self *SRegion) GetZoneSupportedDiskTypes(zoneId string) ([]string, error) {
	dts, err := self.GetDiskTypes()
	if err != nil {
		return nil, errors.Wrap(err, "GetDiskTypes")
	}

	ret := []string{}
	for i := range dts {
		if dts[i].IsAvaliableInZone(zoneId) {
			ret = append(ret, dts[i].Name)
		}
	}

	return ret, nil
}

func (region *SRegion) GetIVMs() ([]cloudprovider.ICloudVM, error) {
	vms, err := region.GetInstances()
	if err != nil {
		return nil, errors.Wrap(err, "GetInstances")
	}
	ret := []cloudprovider.ICloudVM{}
	for i := range vms {
		ret = append(ret, &vms[i])
	}
	return ret, nil
}
