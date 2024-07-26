package qcloud

import (
	"fmt"
	"strings"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/billing"
)

func (self *SHost) _createVMCfel(name, hostname string, imgId string, sysDisk cloudprovider.SDiskInfo, cpu int, memMB int, instanceType string,
	networkId string, ipAddr string, desc string, passwd string,
	diskSizes []cloudprovider.SDiskInfo, publicKey string, secgroupIds []string, userData string, bc *billing.SBillingCycle, projectId string,
	publicIpBw int, publicIpChargeType cloudprovider.TElasticipChargeType,
	tags map[string]string, osType string,
) (string, error) {
	var err error
	keypair := ""
	if len(publicKey) > 0 {
		keypair, err = self.zone.region.syncKeypair(publicKey)
		if err != nil {
			return "", err
		}
	}

	img, err := self.zone.region.GetImage(imgId)
	if err != nil {
		return "", errors.Wrapf(err, "GetImage(%s)", imgId)
	}
	if img.ImageState != ImageStatusNormal && img.ImageState != ImageStatusUsing {
		return "", fmt.Errorf("image %s not ready status is %s", imgId, img.ImageState)
	}

	err = self.zone.validateStorageType(sysDisk.StorageType)
	if err != nil {
		return "", fmt.Errorf("Storage %s not avaiable: %s", sysDisk.StorageType, err)
	}

	disks := make([]SDisk, len(diskSizes)+1)
	disks[0].DiskSize = img.ImageSize
	if sysDisk.SizeGB > 0 && sysDisk.SizeGB > img.ImageSize {
		disks[0].DiskSize = sysDisk.SizeGB
	}

	// 根据实际输入来创建系统盘
	//if disks[0].DiskSize < 50 {
	//	disks[0].DiskSize = 50
	//}

	disks[0].DiskType = strings.ToUpper(sysDisk.StorageType)

	for i, dataDisk := range diskSizes {
		disks[i+1].DiskSize = dataDisk.SizeGB
		err = self.zone.validateStorageType(dataDisk.StorageType)
		if err != nil {
			return "", fmt.Errorf("Storage %s not avaiable: %s", dataDisk.StorageType, err)
		}
		disks[i+1].DiskType = strings.ToUpper(dataDisk.StorageType)
	}

	if len(instanceType) > 0 {
		log.Debugf("Try instancetype : %s", instanceType)
		vmId, err := self.zone.region.CfelCreateInstance(name, hostname, imgId, instanceType, secgroupIds, self.zone.Zone, desc, passwd, disks, networkId, ipAddr, keypair, userData, bc, projectId, publicIpBw, publicIpChargeType, tags, osType)
		if err != nil {
			return "", errors.Wrapf(err, "Failed to create specification %s", instanceType)
		}
		return vmId, nil
	}

	instanceTypes, err := self.zone.region.GetMatchInstanceTypes(cpu, memMB, 0, self.zone.Zone)
	if err != nil {
		return "", err
	}
	if len(instanceTypes) == 0 {
		return "", fmt.Errorf("instance type %dC%dMB not avaiable", cpu, memMB)
	}

	var vmId string
	for _, instType := range instanceTypes {
		instanceTypeId := instType.InstanceType
		log.Debugf("Try instancetype : %s", instanceTypeId)
		vmId, err = self.zone.region.CfelCreateInstance(name, hostname, imgId, instanceTypeId, secgroupIds, self.zone.Zone, desc, passwd, disks, networkId, ipAddr, keypair, userData, bc, projectId, publicIpBw, publicIpChargeType, tags, osType)
		if err != nil {
			log.Errorf("Failed for %s: %s", instanceTypeId, err)
		} else {
			return vmId, nil
		}
	}

	return "", fmt.Errorf("Failed to create, %s", err.Error())
}
