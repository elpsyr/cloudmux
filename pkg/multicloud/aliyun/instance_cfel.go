package aliyun

import (
	"context"
	"fmt"
	"time"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/billing"
	"yunion.io/x/pkg/utils"
)

func (self *SInstance) RebootVM(ctx context.Context) error {
	err := self.host.zone.region.RebootVM(self.InstanceId)
	if err != nil {
		return err
	}
	return cloudprovider.WaitStatus(self, api.VM_RUNNING, 10*time.Second, 300*time.Second) // 5mintues
}

func (self *SInstance) GetIHostId() string {
	return fmt.Sprintf("-%s", self.ZoneId)
}

func (self *SRegion) doRebootVM(instanceId string) error {
	return self.instanceOperation(instanceId, "RebootInstance", nil)
}

func (self *SRegion) RebootVM(instanceId string) error {
	status, err := self.GetInstanceStatus(instanceId)
	if err != nil {
		log.Errorf("Fail to get instance status on StartVM: %s", err)
		return err
	}
	if status != InstanceStatusRunning {
		log.Errorf("RebootVM: vm status is %s expect %s", status, InstanceStatusRunning)
		return cloudprovider.ErrInvalidStatus
	}
	return self.doRebootVM(instanceId)
	// if err != nil {
	//	return err
	// }
	// return self.waitInstanceStatus(instanceId, InstanceStatusRunning, time.Second*5, time.Second*180) // 3 minutes to timeout
}

func (self *SInstance) GetMonitorData(start, end string) ([]cloudprovider.ICfelMonitorData, error) {
	data, err := self.host.zone.region.DescribeInstanceMonitorData(self.InstanceId, start, end, "")
	if err != nil {
		return nil, errors.Wrap(err, "DescribeInstanceMonitorData")
	}
	// 将 []MonitorDataItem 转换成 []cloudprovider.MonitorData
	var providerData []cloudprovider.ICfelMonitorData
	for _, item := range data {
		providerData = append(providerData, cloudprovider.ICfelMonitorData(item))
	}
	return providerData, nil
}

func (self *SRegion) DescribeInstanceMonitorData(instanceId, startTime, endTime, period string) ([]MonitorDataItem, error) {
	params := map[string]string{
		"InstanceId": instanceId,
		"StartTime":  startTime,
		"EndTime":    endTime,
		"Period":     period,
	}
	if period == "" {
		params["Period"] = "600" // 默认值：60
	}
	resp, err := self.ecsRequest("DescribeInstanceMonitorData", params)
	if err != nil {
		return nil, errors.Wrapf(err, "DescribeInstanceMonitorData")
	}
	ret := DescribeInstanceMonitorDataResponse{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, errors.Wrapf(err, "resp.Unmarshal")
	}
	return ret.MonitorData.InstanceMonitorData, nil
}

// DescribeInstanceMonitorDataResponse
// DescribeInstanceMonitorData 接口返回数据结构
type DescribeInstanceMonitorDataResponse struct {
	RequestId   string `json:"RequestId"`
	MonitorData struct {
		InstanceMonitorData []MonitorDataItem `json:"InstanceMonitorData"`
	} `json:"MonitorData"`
}

type MonitorDataItem struct {
	IOPSRead          float64 `json:"IOPSRead,omitempty"`
	IntranetBandwidth float64 `json:"IntranetBandwidth"`
	IOPSWrite         float64 `json:"IOPSWrite,omitempty"`
	InstanceId        string  `json:"InstanceId"`
	IntranetTX        float64 `json:"IntranetTX"`
	CPU               float64 `json:"CPU"`
	BPSRead           float64 `json:"BPSRead,omitempty"`
	IntranetRX        float64 `json:"IntranetRX"`
	TimeStamp         string  `json:"TimeStamp"`
	InternetBandwidth float64 `json:"InternetBandwidth"`
	InternetTX        float64 `json:"InternetTX"`
	InternetRX        float64 `json:"InternetRX"`
	BPSWrite          float64 `json:"BPSWrite,omitempty"`
}

var _ cloudprovider.ICfelMonitorData = (*MonitorDataItem)(nil)

func (m MonitorDataItem) GetBPSRead() float64 {
	return m.IOPSRead
}

func (m MonitorDataItem) GetInternetTX() float64 {
	return m.InternetTX
}

func (m MonitorDataItem) GetCPU() float64 {
	return m.CPU
}

func (m MonitorDataItem) GetMem() float64 {
	return 0
}

func (m MonitorDataItem) GetDisk() float64 {
	return 0
}

func (m MonitorDataItem) GetIOPSWrite() float64 {
	return m.IOPSWrite
}

func (m MonitorDataItem) GetIntranetTX() float64 {
	return m.IntranetTX
}

func (m MonitorDataItem) GetInstanceId() string {
	return m.InstanceId
}

func (m MonitorDataItem) GetBPSWrite() float64 {
	return m.BPSWrite
}

func (m MonitorDataItem) GetIOPSRead() float64 {
	return m.IOPSRead
}

func (m MonitorDataItem) GetInternetBandwidth() float64 {
	return m.InternetBandwidth
}

func (m MonitorDataItem) GetInternetRX() float64 {
	return m.InternetRX
}

func (m MonitorDataItem) GetTimeStamp() string {
	return m.TimeStamp
}

func (m MonitorDataItem) GetIntranetRX() float64 {
	return m.IntranetRX
}

func (m MonitorDataItem) GetIntranetBandwidth() float64 {
	return m.IntranetBandwidth
}

func (self *SRegion) CreateInstanceCfel(name, hostname string, imageId string, instanceType string, securityGroupIds []string,
	zoneId string, desc string, passwd string, disks []SDisk, vSwitchId string, ipAddr string,
	keypair string, userData string, bc *billing.SBillingCycle, projectId, osType string,
	tags map[string]string, publicIp cloudprovider.SPublicIpInfo,
) (string, error) {
	params := make(map[string]string)
	params["RegionId"] = self.RegionId
	params["ImageId"] = imageId
	params["InstanceType"] = instanceType
	for _, id := range securityGroupIds {
		params["SecurityGroupId"] = id
	}
	params["ZoneId"] = zoneId
	params["InstanceName"] = name
	if len(hostname) > 0 {
		params["HostName"] = hostname
	}
	params["Description"] = desc
	params["InternetChargeType"] = "PayByTraffic"
	if publicIp.PublicIpBw > 0 {
		params["InternetMaxBandwidthOut"] = fmt.Sprintf("%d", publicIp.PublicIpBw)
		params["InternetMaxBandwidthIn"] = "200"
	}
	if publicIp.PublicIpChargeType == cloudprovider.ElasticipChargeTypeByBandwidth {
		params["InternetChargeType"] = "PayByBandwidth"
	}
	if len(passwd) > 0 {
		params["Password"] = passwd
	} else {
		params["PasswordInherit"] = "True"
	}

	if len(projectId) > 0 {
		params["ResourceGroupId"] = projectId
	}
	//{"Code":"InvalidSystemDiskCategory.ValueNotSupported","HostId":"ecs.aliyuncs.com","Message":"The specified parameter 'SystemDisk.Category' is not support IoOptimized Instance. Valid Values: cloud_efficiency;cloud_ssd. ","RequestId":"9C9A4E99-5196-42A2-80B6-4762F8F75C90"}
	params["IoOptimized"] = "optimized"
	for i, d := range disks {
		if i == 0 {
			params["SystemDisk.Category"] = d.Category
			if d.Category == api.STORAGE_CLOUD_ESSD_PL0 {
				params["SystemDisk.Category"] = api.STORAGE_CLOUD_ESSD
				params["SystemDisk.PerformanceLevel"] = "PL0"
			}
			if d.Category == api.STORAGE_CLOUD_ESSD_PL2 {
				params["SystemDisk.Category"] = api.STORAGE_CLOUD_ESSD
				params["SystemDisk.PerformanceLevel"] = "PL2"
			}
			if d.Category == api.STORAGE_CLOUD_ESSD_PL3 {
				params["SystemDisk.Category"] = api.STORAGE_CLOUD_ESSD
				params["SystemDisk.PerformanceLevel"] = "PL3"
			}
			if d.Category == api.STORAGE_CLOUD_AUTO {
				params["SystemDisk.BurstingEnabled"] = "true"
			}
			params["SystemDisk.Size"] = fmt.Sprintf("%d", d.Size)
			params["SystemDisk.DiskName"] = d.GetName()
			params["SystemDisk.Description"] = d.Description
		} else {
			params[fmt.Sprintf("DataDisk.%d.Size", i)] = fmt.Sprintf("%d", d.Size)
			params[fmt.Sprintf("DataDisk.%d.Category", i)] = d.Category
			if d.Category == api.STORAGE_CLOUD_ESSD_PL0 {
				params[fmt.Sprintf("DataDisk.%d.Category", i)] = api.STORAGE_CLOUD_ESSD
				params[fmt.Sprintf("DataDisk.%d..PerformanceLevel", i)] = "PL0"
			}
			if d.Category == api.STORAGE_CLOUD_ESSD_PL2 {
				params[fmt.Sprintf("DataDisk.%d.Category", i)] = api.STORAGE_CLOUD_ESSD
				params[fmt.Sprintf("DataDisk.%d..PerformanceLevel", i)] = "PL2"
			}
			if d.Category == api.STORAGE_CLOUD_ESSD_PL3 {
				params[fmt.Sprintf("DataDisk.%d.Category", i)] = api.STORAGE_CLOUD_ESSD
				params[fmt.Sprintf("DataDisk.%d..PerformanceLevel", i)] = "PL3"
			}
			if d.Category == api.STORAGE_CLOUD_AUTO {
				params[fmt.Sprintf("DataDisk.%d.BurstingEnabled", i)] = "true"
			}
			params[fmt.Sprintf("DataDisk.%d.DiskName", i)] = d.GetName()
			params[fmt.Sprintf("DataDisk.%d.Description", i)] = d.Description
			params[fmt.Sprintf("DataDisk.%d.Encrypted", i)] = "false"
		}
	}
	params["VSwitchId"] = vSwitchId
	params["PrivateIpAddress"] = ipAddr

	if len(keypair) > 0 {
		params["KeyPairName"] = keypair
	}

	if len(userData) > 0 {
		params["UserData"] = userData
	}

	if len(tags) > 0 {
		tagIdx := 1
		for k, v := range tags {
			params[fmt.Sprintf("Tag.%d.Key", tagIdx)] = k
			params[fmt.Sprintf("Tag.%d.Value", tagIdx)] = v
			tagIdx += 1
		}
	}

	if bc != nil {
		params["InstanceChargeType"] = "PrePaid"
		err := billingCycle2Params(bc, params)
		if err != nil {
			return "", err
		}
		if bc.AutoRenew {
			params["AutoRenew"] = "true"
			params["AutoRenewPeriod"] = "1"
		} else {
			params["AutoRenew"] = "False"
		}
	} else {
		params["InstanceChargeType"] = "PostPaid"
		params["SpotStrategy"] = "NoSpot"
	}

	ct, ok := tags[cloudprovider.InstanceChargeTypeTag]
	if ok && ct == cloudprovider.InstanceChargeTypeSpotPaid {
		params["InstanceChargeType"] = "PostPaid"
		params["SpotStrategy"] = "SpotAsPriceGo"
	}

	params["ClientToken"] = utils.GenRequestId(20)

	resp, err := self.ecsRequest("RunInstances", params)
	if err != nil {
		return "", errors.Wrapf(err, "RunInstances")
	}
	ids := []string{}
	err = resp.Unmarshal(&ids, "InstanceIdSets", "InstanceIdSet")
	if err != nil {
		return "", errors.Wrapf(err, "Unmarshal")
	}
	for _, id := range ids {
		err = cloudprovider.Wait(time.Second*3, time.Minute, func() (bool, error) {
			_, err := self.GetInstance(id)
			if err != nil {
				if errors.Cause(err) == cloudprovider.ErrNotFound {
					return false, nil
				}
				return false, err
			}
			return true, nil
		})
		if err != nil {
			return "", errors.Wrapf(cloudprovider.ErrNotFound, "after vm %s created", id)
		}
		return id, nil
	}
	return "", errors.Wrapf(cloudprovider.ErrNotFound, "after created")
}
