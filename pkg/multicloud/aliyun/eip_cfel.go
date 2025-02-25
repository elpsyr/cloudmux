package aliyun

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/utils"
)

func (region *SRegion) AllocateEIP(opts *cloudprovider.SEip) (*SEipAddress, error) {
	params := make(map[string]string)
	if len(opts.Name) > 0 {
		params["Name"] = opts.Name
	}
	params["RegionId"] = region.RegionId
	params["Bandwidth"] = fmt.Sprintf("%d", opts.BandwidthMbps)
	if strings.ToLower(opts.ChargeType) == "postpaid" { // 按量计费
		switch opts.InternetChargeType {
		case api.EIP_CHARGE_TYPE_BY_TRAFFIC:
			params["InternetChargeType"] = string(InternetChargeByTraffic)
		case api.EIP_CHARGE_TYPE_BY_BANDWIDTH:
			params["InternetChargeType"] = string(InternetChargeByBandwidth)

		}
	} else { // 包年包月
		var pricingCycle = "Month"
		if strings.ToLower(opts.PricingCycle) == "year" {
			pricingCycle = "Year"
		}
		params["PricingCycle"] = pricingCycle
		params["Period"] = strconv.Itoa(opts.Period)
		params["AutoPay"] = "true" // 包年包月开启自动付费
	}
	params["InstanceChargeType"] = opts.ChargeType
	params["ClientToken"] = utils.GenRequestId(20)
	if len(opts.ProjectId) > 0 {
		params["ResourceGroupId"] = opts.ProjectId
	}
	params["ISP"] = "BGP"
	if opts.BGPType == "BGP_PRO" {
		params["ISP"] = "BGP_PRO"
	}

	body, err := region.vpcRequest("AllocateEipAddress", params)
	if err != nil {
		return nil, errors.Wrapf(err, "AllocateEipAddress")
	}

	eipId, err := body.GetString("AllocationId")
	if err != nil {
		return nil, errors.Wrapf(err, "get AllocationId after created")
	}
	var eip *SEipAddress
	var i = 0
	for i < 4 { // 包年包月会比较慢
		i ++
		eip, err = region.GetEip(eipId)
		if err != nil {
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	
	if err != nil {
		return nil, err
	}
	if len(opts.Tags) == 0 {
		return eip, nil
	}
	cloudprovider.WaitStatus(eip, api.EIP_STATUS_READY, time.Second*5, time.Minute*1)
	eip.SetTags(opts.Tags, false)
	return eip, nil
}
