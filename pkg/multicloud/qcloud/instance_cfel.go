package qcloud

import (
	"context"
	"fmt"
	"strings"
	"time"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/billing"
	"yunion.io/x/pkg/utils"
)

func (self *SInstance) GetIHostId() string {
	return fmt.Sprintf("-%s", self.Placement.Zone)
}

func (self *SRegion) CfelCreateInstance(name, hostname string, imageId string, instanceType string, securityGroupIds []string,
	zoneId string, desc string, passwd string, disks []SDisk, networkId string, ipAddr string,
	keypair string, userData string, bc *billing.SBillingCycle, projectId string,
	publicIpBw int, publicIpChargeType cloudprovider.TElasticipChargeType,
	tags map[string]string, osType string,
) (string, error) {
	params := make(map[string]string)
	params["Region"] = self.Region
	params["ImageId"] = imageId
	params["InstanceType"] = instanceType
	for i, id := range securityGroupIds {
		params[fmt.Sprintf("SecurityGroupIds.%d", i)] = id
	}
	params["Placement.Zone"] = zoneId
	if len(projectId) > 0 {
		params["Placement.ProjectId"] = projectId
	}
	params["InstanceName"] = name
	if len(hostname) > 0 {
		params["HostName"] = hostname
	}

	bandwidth := publicIpBw
	if publicIpChargeType != "" && publicIpBw == 0 {
		bandwidth = 100
		if bc != nil {
			bandwidth = 200
		}
	}

	internetChargeType := "TRAFFIC_POSTPAID_BY_HOUR"
	if publicIpChargeType == cloudprovider.ElasticipChargeTypeByBandwidth {
		internetChargeType = "BANDWIDTH_POSTPAID_BY_HOUR"
	}
	pkgs, _, err := self.GetBandwidthPackages([]string{}, 0, 50)
	if err != nil {
		return "", errors.Wrapf(err, "GetBandwidthPackages")
	}
	if len(pkgs) > 0 {
		bandwidth = 65535   // unlimited bandwidth
		if publicIpBw > 0 { // 若用户指定带宽则限制带宽大小
			bandwidth = publicIpBw
		}
		internetChargeType = "BANDWIDTH_PACKAGE"
		pkgId := pkgs[0].BandwidthPackageId
		for _, pkg := range pkgs {
			if len(pkg.ResourceSet) < 100 {
				pkgId = pkg.BandwidthPackageId
				break
			}
		}
		params["InternetAccessible.BandwidthPackageId"] = pkgId
	}

	params["InternetAccessible.InternetChargeType"] = internetChargeType
	params["InternetAccessible.InternetMaxBandwidthOut"] = fmt.Sprintf("%d", bandwidth)
	params["InternetAccessible.PublicIpAssigned"] = "TRUE"
	if publicIpBw == 0 {
		params["InternetAccessible.PublicIpAssigned"] = "FALSE"
	}
	if len(keypair) > 0 {
		params["LoginSettings.KeyIds.0"] = keypair
	} else if len(passwd) > 0 {
		params["LoginSettings.Password"] = passwd
	} else {
		params["LoginSettings.KeepImageLogin"] = "TRUE"
	}
	if len(userData) > 0 {
		params["UserData"] = userData
	}

	if bc != nil {
		params["InstanceChargeType"] = "PREPAID"
		params["InstanceChargePrepaid.Period"] = fmt.Sprintf("%d", bc.GetMonths())
		if bc.AutoRenew {
			params["InstanceChargePrepaid.RenewFlag"] = "NOTIFY_AND_AUTO_RENEW"
		} else {
			params["InstanceChargePrepaid.RenewFlag"] = "NOTIFY_AND_MANUAL_RENEW"
		}
	} else {
		params["InstanceChargeType"] = "POSTPAID_BY_HOUR"
	}
	ct, ok := tags[cloudprovider.InstanceChargeTypeTag]
	if ok && ct == cloudprovider.InstanceChargeTypeSpotPaid {
		params["InstanceChargeType"] = "SPOTPAID"
	}

	// tags
	if len(tags) > 0 {
		params["TagSpecification.0.ResourceType"] = "instance"
		tagIdx := 0
		for k, v := range tags {
			params[fmt.Sprintf("TagSpecification.0.Tags.%d.Key", tagIdx)] = k
			params[fmt.Sprintf("TagSpecification.0.Tags.%d.Value", tagIdx)] = v
			tagIdx += 1
		}
	}

	//params["IoOptimized"] = "optimized"
	for i, d := range disks {
		if i == 0 {
			params["SystemDisk.DiskType"] = d.DiskType
			params["SystemDisk.DiskSize"] = fmt.Sprintf("%d", d.DiskSize)
		} else {
			params[fmt.Sprintf("DataDisks.%d.DiskSize", i-1)] = fmt.Sprintf("%d", d.DiskSize)
			params[fmt.Sprintf("DataDisks.%d.DiskType", i-1)] = d.DiskType
		}
	}
	network, err := self.GetNetwork(networkId)
	if err != nil {
		return "", errors.Wrapf(err, "GetNetwork(%s)", networkId)
	}
	params["VirtualPrivateCloud.SubnetId"] = networkId
	params["VirtualPrivateCloud.VpcId"] = network.VpcId
	if len(ipAddr) > 0 {
		params["VirtualPrivateCloud.PrivateIpAddresses.0"] = ipAddr
	}

	var body jsonutils.JSONObject
	instanceIdSet := []string{}
	err = cloudprovider.Wait(time.Second*10, time.Minute, func() (bool, error) {
		params["ClientToken"] = utils.GenRequestId(20)
		body, err = self.cvmRequest("RunInstances", params, true)
		if err != nil {
			if strings.Contains(err.Error(), "Code=InvalidPermission") { // 带宽上移用户未指定公网ip时不能设置带宽
				delete(params, "InternetAccessible.InternetChargeType")
				delete(params, "InternetAccessible.InternetMaxBandwidthOut")
				return false, nil
			}
			if strings.Contains(err.Error(), "UnsupportedOperation.BandwidthPackageIdNotSupported") ||
				(strings.Contains(err.Error(), "Code=InvalidParameterCombination") && strings.Contains(err.Error(), "InternetAccessible.BandwidthPackageId")) {
				delete(params, "InternetAccessible.BandwidthPackageId")
				return false, nil
			}
			return false, errors.Wrapf(err, "RunInstances")
		}
		return true, nil
	})
	if err != nil {
		return "", errors.Wrap(err, "RunInstances")
	}
	err = body.Unmarshal(&instanceIdSet, "InstanceIdSet")
	if err == nil && len(instanceIdSet) > 0 {
		return instanceIdSet[0], nil
	}
	return "", fmt.Errorf("Failed to create instance")
}

func (self *SInstance) RebootVM(ctx context.Context) error {
	err := self.host.zone.region.RebootVM(self.InstanceId)
	if err != nil {
		return err
	}
	return cloudprovider.WaitStatus(self, api.VM_RUNNING, 10*time.Second, 300*time.Second) // 5mintues
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

}

func (self *SRegion) doRebootVM(instanceId string) error {
	return self.instanceOperation(instanceId, "RebootInstances", nil, true)
}
