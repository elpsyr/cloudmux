package qcloud

import (
	"fmt"
	"strconv"
	"strings"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

func (region *SRegion) AllocateEIP(opts *cloudprovider.SEip) (*SEipAddress, error) {
	params := make(map[string]string)
	params["AddressName"] = opts.Name
	if len(opts.Name) > 20 {
		params["AddressName"] = opts.Name[:20]
	}
	if opts.BandwidthMbps > 0 {
		params["InternetMaxBandwidthOut"] = fmt.Sprintf("%d", opts.BandwidthMbps)
	}

	if strings.ToLower(opts.ChargeType) == "postpaid" { // 按量计费
		switch opts.InternetChargeType {
		case api.EIP_CHARGE_TYPE_BY_TRAFFIC:
			params["InternetChargeType"] = "TRAFFIC_POSTPAID_BY_HOUR"
		case api.EIP_CHARGE_TYPE_BY_BANDWIDTH:
			params["InternetChargeType"] = "BANDWIDTH_POSTPAID_BY_HOUR"
		}
	} else { // 包年包月
		
		var period = opts.Period
		if strings.ToLower(opts.PricingCycle) == "year" {
			period *= 12
		}
		params["InternetChargeType"] = "BANDWIDTH_PREPAID_BY_MONTH"
		params["AddressChargePrepaid.Period"] = strconv.Itoa(period)
		params["AddressChargePrepaid.AutoRenewFlag"] = "0" // 自动续费标志。0表示手动续费，1表示自动续费，2表示到期不续费。默认缺省为0即手动续费
	}
	
	idx := 0
	for k, v := range opts.Tags {
		params[fmt.Sprintf("Tags.%d.Key", idx)] = k
		params[fmt.Sprintf("Tags.%d.Value", idx)] = v
		idx++
	}

	addRessSet := []string{}
	body, err := region.vpcRequest("AllocateAddresses", params)
	if err != nil {
		return nil, errors.Wrapf(err, "AllocateAddresses")
	}
	err = body.Unmarshal(&addRessSet, "AddressSet")
	if err != nil {
		return nil, errors.Wrapf(err, "resp.Unmarshal")
	}
	return region.GetEip(addRessSet[0])
}
